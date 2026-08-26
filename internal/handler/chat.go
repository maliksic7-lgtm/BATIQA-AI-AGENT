package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"batiqa-ai/internal/model"
	"batiqa-ai/internal/repository"
	"batiqa-ai/internal/service/ai"
	ticketservice "batiqa-ai/internal/service/ticket"
)

// ChatHandler handles POST /api/chat per AI CHAT.MD and USER_FLOW.md
// Flow: Guest Message -> AI -> Structured AI Result -> Backend Validation -> Ticket Service -> Repository -> DB
// For info intents, it retrieves verified hotel information from DB per AI KNOWLEDGE SOURCE.md
// For recommendation intents, it retrieves verified local data from DB per RECOMMENDATION RULES.md
type ChatHandler struct {
	ai         *ai.Service
	tickets    *ticketservice.Service
	convs      *repository.ConversationRepository
	guests     *repository.GuestRepository
	hotelInfos *repository.HotelInfoRepository
	recs       *repository.RecommendationRepository
}

func NewChatHandler(aiSvc *ai.Service, ticketSvc *ticketservice.Service, convRepo *repository.ConversationRepository, guestRepo *repository.GuestRepository) *ChatHandler {
	return &ChatHandler{ai: aiSvc, tickets: ticketSvc, convs: convRepo, guests: guestRepo}
}

func NewChatHandlerWithHotel(aiSvc *ai.Service, ticketSvc *ticketservice.Service, convRepo *repository.ConversationRepository, guestRepo *repository.GuestRepository, hotelRepo *repository.HotelInfoRepository) *ChatHandler {
	return &ChatHandler{ai: aiSvc, tickets: ticketSvc, convs: convRepo, guests: guestRepo, hotelInfos: hotelRepo}
}

func NewChatHandlerFull(aiSvc *ai.Service, ticketSvc *ticketservice.Service, convRepo *repository.ConversationRepository, guestRepo *repository.GuestRepository, hotelRepo *repository.HotelInfoRepository, recRepo *repository.RecommendationRepository) *ChatHandler {
	return &ChatHandler{ai: aiSvc, tickets: ticketSvc, convs: convRepo, guests: guestRepo, hotelInfos: hotelRepo, recs: recRepo}
}

// ChatRequest per AI CHAT.MD
type ChatRequest struct {
	SessionID  string `json:"session_id"`
	RoomNumber string `json:"room_number"`
	Message    string `json:"message"`
}

// ChatResponse per AI CHAT.MD
type ChatResponse struct {
	Message        string  `json:"message"`
	Intent         string  `json:"intent"`
	RequiresTicket bool    `json:"requires_ticket"`
	TicketID       *string `json:"ticket_id"`
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}
	req.SessionID = strings.TrimSpace(req.SessionID)
	req.RoomNumber = strings.TrimSpace(req.RoomNumber)
	req.Message = strings.TrimSpace(req.Message)

	if req.SessionID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "session_id is required")
		return
	}
	if req.Message == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "message is required")
		return
	}
	if len(req.Message) > 2000 {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "message too long (max 2000)")
		return
	}

	// Persist guest session (upsert) - room not invented.
	// Non-fatal: guest table is auxiliary; never block the chat flow on it.
	var roomPtr *string
	if req.RoomNumber != "" {
		roomPtr = &req.RoomNumber
	}
	lang := ai.DetectLanguage(req.Message)
	if _, err := h.guests.Upsert(req.SessionID, roomPtr, lang); err != nil {
		fmt.Printf("guest upsert failed (non-fatal): %v\n", err)
	}

	// Store user conversation (best effort)
	_ = h.convs.Create(&model.Conversation{
		SessionID: req.SessionID,
		Role:      model.RoleUser,
		Message:   req.Message,
		Intent:    nil,
	})

	// Call AI service (with validation, fallback, no crash)
	aiReq := ai.Request{
		SessionID:  req.SessionID,
		RoomNumber: req.RoomNumber,
		Message:    req.Message,
	}
	aiResult, err := h.ai.Process(r.Context(), aiReq)
	if err != nil || aiResult == nil {
		// Defensive: Process guarantees non-nil result, but guard anyway per ERROR FLOW.md
		aiResult = &ai.AIResult{
			Intent:   ai.IntentUnknown,
			Language: ai.LangID,
			Entities: map[string]interface{}{},
			Action:   ai.Action{Type: ai.ActionClarify},
			Response: "Maaf, layanan AI sedang mengalami gangguan. Silakan coba kembali beberapa saat lagi.",
		}
	}

	// For hotel info intents, retrieve verified data from DB per AI KNOWLEDGE SOURCE.md (do not invent)
	if h.hotelInfos != nil && isInfoIntent(aiResult.Intent) {
		if infoResp := h.getHotelInfoResponse(aiResult.Intent, aiResult.Language); infoResp != "" {
			aiResult.Response = infoResp
		} else {
			// If no data, use fallback per spec
			if aiResult.Language == ai.LangEN {
				aiResult.Response = "Sorry, I don't have that information yet. Please contact Front Office."
			} else {
				aiResult.Response = "Maaf, saya belum memiliki informasi tersebut. Silakan hubungi Front Office."
			}
		}
	}

	// For recommendation intents, retrieve verified local data from DB per RECOMMENDATION RULES.md
	if h.recs != nil && isRecommendationIntent(aiResult.Intent) && aiResult.Action.Type == ai.ActionAnswer {
		h.enrichRecommendationResponse(aiResult)
	}

	// Store assistant conversation (best effort)
	_ = h.convs.Create(&model.Conversation{
		SessionID: req.SessionID,
		Role:      model.RoleAssistant,
		Message:   aiResult.Response,
		Intent:    &aiResult.Intent,
	})

	// If AI requires ticket, attempt to create via Ticket Service (backend validation)
	var ticketID *string
	requiresTicket := aiResult.RequiresTicket

	if requiresTicket {
		// Check missing room_number: if not in entities and no provided room, do not create, ask for room
		if _, ok := aiResult.Entities["room_number"]; !ok && req.RoomNumber == "" {
			// Missing room flow per MISSING ROOM NUMBER.md
			requiresTicket = true
			ticketID = nil
		} else {
			t, err := h.tickets.CreateFromAI(aiResult, req.Message)
			if err != nil {
				switch {
				case errors.Is(err, ticketservice.ErrRoomRequired):
					requiresTicket = true
					ticketID = nil
					if !strings.Contains(strings.ToLower(aiResult.Response), "kamar") && !strings.Contains(strings.ToLower(aiResult.Response), "room") {
						if aiResult.Language == ai.LangEN {
							aiResult.Response = "Sure, I can help. Could you please tell me your room number?"
						} else {
							aiResult.Response = "Baik, saya bisa membantu. Boleh saya tahu nomor kamar Anda?"
						}
					}
				case ticketservice.IsValidationError(err):
					// Validation error -> do not create ticket
					requiresTicket = false
					ticketID = nil
				default:
					// Internal DB error -> keep requires_ticket true but no ticket, inform guest of temporary failure
					requiresTicket = true
					ticketID = nil
					if aiResult.Language == ai.LangEN {
						aiResult.Response = "Sorry, I couldn't create your ticket due to a temporary issue. Please try again shortly or contact Front Office."
					} else {
						aiResult.Response = "Maaf, saya tidak bisa membuat tiket karena gangguan sementara. Silakan coba lagi atau hubungi Front Office."
					}
				}
			} else {
				s := t.TicketNumber
				ticketID = &s
			}
		}
	}

	resp := ChatResponse{
		Message:        aiResult.Response,
		Intent:         aiResult.Intent,
		RequiresTicket: requiresTicket,
		TicketID:       ticketID,
	}
	WriteJSON(w, http.StatusOK, resp)
}

func isInfoIntent(intent string) bool {
	switch intent {
	case ai.IntentBreakfastInformation, ai.IntentWifiInformation, ai.IntentCheckinInformation,
		ai.IntentCheckoutInformation, ai.IntentFacilityInformation, ai.IntentHotelInformation,
		ai.IntentRestaurantInformation, ai.IntentRoomInformation, ai.IntentHotelPolicy:
		return true
	default:
		return false
	}
}

func isRecommendationIntent(intent string) bool {
	switch intent {
	case ai.IntentRestaurantRecommendation, ai.IntentCafeRecommendation, ai.IntentTourismRecommendation,
		ai.IntentShoppingRecommendation, ai.IntentATMRequest, ai.IntentTransportationRequest:
		return true
	default:
		return false
	}
}

func recommendationCategory(intent string) string {
	switch intent {
	case ai.IntentRestaurantRecommendation:
		return "restaurant"
	case ai.IntentCafeRecommendation:
		return "cafe"
	case ai.IntentTourismRecommendation:
		return "tourism"
	case ai.IntentShoppingRecommendation:
		return "shopping"
	case ai.IntentATMRequest:
		return "atm"
	case ai.IntentTransportationRequest:
		return "transportation"
	default:
		return ""
	}
}

// enrichRecommendationResponse replaces the generic template with verified data from the
// recommendations table (top 3 nearest). Never invents entries per RECOMMENDATION RULES.md.
func (h *ChatHandler) enrichRecommendationResponse(res *ai.AIResult) {
	cat := recommendationCategory(res.Intent)
	if cat == "" {
		return
	}
	var maxPrice *int
	if v, ok := res.Entities["budget"]; ok {
		if n := toInt(v); n > 0 {
			maxPrice = &n
		}
	}
	items, err := h.recs.ListActive(&cat, maxPrice)
	if err != nil {
		fmt.Printf("recommendation query failed (non-fatal): %v\n", err)
		return
	}
	if len(items) == 0 {
		if res.Language == ai.LangEN {
			res.Response = "Sorry, I don't have recommendations for that yet. Please contact Front Office."
		} else {
			res.Response = "Maaf, saya belum punya rekomendasi untuk itu. Silakan hubungi Front Office."
		}
		return
	}
	if len(items) > 3 {
		items = items[:3]
	}
	var b strings.Builder
	if res.Language == ai.LangEN {
		b.WriteString("Here are nearby recommendations:")
	} else {
		b.WriteString("Berikut rekomendasi terdekat untuk Anda:")
	}
	for _, it := range items {
		b.WriteString("\n• ")
		b.WriteString(it.Name)
		if it.Description != nil && *it.Description != "" {
			b.WriteString(" — ")
			b.WriteString(*it.Description)
		}
		details := make([]string, 0, 2)
		if price := formatPriceRange(it.PriceMin, it.PriceMax); price != "" {
			details = append(details, price)
		}
		if it.DistanceKm != nil {
			details = append(details, fmt.Sprintf("~%s km", trimFloat(*it.DistanceKm)))
		}
		if len(details) > 0 {
			b.WriteString(" (")
			b.WriteString(strings.Join(details, ", "))
			b.WriteString(")")
		}
	}
	res.Response = b.String()
}

func toInt(v interface{}) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case string:
		var out int
		_, err := fmt.Sscanf(strings.TrimSpace(n), "%d", &out)
		if err != nil {
			return 0
		}
		return out
	default:
		return 0
	}
}

func formatPriceRange(pMin, pMax *int) string {
	switch {
	case pMin != nil && pMax != nil && *pMax > 0:
		return fmt.Sprintf("Rp%s–Rp%s", thousands(*pMin), thousands(*pMax))
	case pMax != nil && *pMax > 0:
		return fmt.Sprintf("≤ Rp%s", thousands(*pMax))
	case pMin != nil && *pMin > 0:
		return fmt.Sprintf("≥ Rp%s", thousands(*pMin))
	default:
		return ""
	}
}

func thousands(n int) string {
	s := strconv.Itoa(n)
	var b strings.Builder
	for i, d := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(d)
	}
	return b.String()
}

func trimFloat(f float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", f), "0"), ".")
}

func (h *ChatHandler) getHotelInfoResponse(intent, lang string) string {
	categoryMap := map[string]string{
		ai.IntentBreakfastInformation:  "BREAKFAST",
		ai.IntentWifiInformation:       "WIFI",
		ai.IntentCheckinInformation:    "CHECKIN",
		ai.IntentCheckoutInformation:   "CHECKOUT",
		ai.IntentFacilityInformation:   "ROOM",
		ai.IntentHotelInformation:      "ROOM",
		ai.IntentRestaurantInformation: "RESTAURANT",
		ai.IntentRoomInformation:       "ROOM",
		ai.IntentHotelPolicy:           "POLICY",
	}
	cat, ok := categoryMap[intent]
	if !ok {
		return ""
	}
	items, err := h.hotelInfos.ListActive(&cat)
	if err != nil || len(items) == 0 {
		return ""
	}
	// Return first item's content (verified DB data)
	// For breakfast, prefer schedule
	for _, it := range items {
		if cat == "BREAKFAST" && strings.Contains(strings.ToLower(it.Title), "schedule") {
			return it.Content
		}
	}
	return items[0].Content
}
