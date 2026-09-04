package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

// ChatRequest per AI CHAT.MD. Token is optional here when the router-level
// guest middleware already validated the header/query token.
type ChatRequest struct {
	SessionID  string `json:"session_id"`
	RoomNumber string `json:"room_number"`
	Message    string `json:"message"`
	Token      string `json:"token,omitempty"`
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

	// Verified QR-token room wins over anything the client claims; guests
	// cannot create tickets or chats for a different room.
	resolvedRoom := GuestRoom(r)
	if resolvedRoom == "" {
		resolvedRoom = req.RoomNumber // legacy/dev path without token
	}

	// Persist guest session (upsert) - room not invented.
	// Non-fatal: guest data is auxiliary; never block the chat flow on it.
	var roomPtr *string
	if resolvedRoom != "" {
		roomPtr = &resolvedRoom
	}
	if _, err := h.guests.Upsert(req.SessionID, roomPtr, ai.DetectLanguage(req.Message)); err != nil {
		fmt.Printf("guest upsert failed (non-fatal): %v\n", err)
	}

	// Save user message AFTER loading history so the AI request contains prior turns.
	aiReq := h.newAIRequest(req.SessionID, resolvedRoom, req.Message)
	_ = h.convs.Create(&model.Conversation{
		SessionID: req.SessionID,
		Role:      model.RoleUser,
		Message:   req.Message,
		Intent:    nil,
	})

	// Call AI service (with validation, fallback, no crash).
	// Verified hotel facts + recommendations are injected so the LLM can compose
	// natural answers from data instead of templated text per AI KNOWLEDGE SOURCE.md.
	aiResult, err := h.ai.Process(r.Context(), *aiReq)
	if err != nil {
		fmt.Printf("AI fallback engaged: %v\n", err)
	}
	if err != nil || aiResult == nil {
		// Defensive: Process guarantees non-nil result, but guard anyway per ERROR FLOW.md
		fb := aiFallbackResult()
		aiResult = &fb
	}

	// Rule-based provider has no access to injected facts, so answers are composed
	// from DB templates here. With Gemini active we trust its fact-grounded phrasing.
	h.applyRuleBasedOverrides(aiResult)

	// Store assistant conversation (best effort)
	_ = h.convs.Create(&model.Conversation{
		SessionID: req.SessionID,
		Role:      model.RoleAssistant,
		Message:   aiResult.Response,
		Intent:    &aiResult.Intent,
	})

	ticketID := h.maybeCreateTicket(r, aiResult, req.Message, resolvedRoom)
	requiresTicket := aiResult.RequiresTicket

	resp := ChatResponse{
		Message:        aiResult.Response,
		Intent:         aiResult.Intent,
		RequiresTicket: requiresTicket,
		TicketID:       ticketID,
	}
	WriteJSON(w, http.StatusOK, resp)
}

// newAIRequest assembles the AI request with conversation history and verified facts.
func (h *ChatHandler) newAIRequest(sessionID, room, message string) *ai.Request {
	return &ai.Request{
		SessionID:  sessionID,
		RoomNumber: room,
		Message:    message,
		History:    loadHistory(h.convs, sessionID, 10),
		Facts:      h.verifiedFacts(),
	}
}

// aiFallbackResult is the controlled ERROR FLOW response.
func aiFallbackResult() ai.AIResult {
	return ai.AIResult{
		Intent:   ai.IntentUnknown,
		Language: ai.LangID,
		Entities: map[string]interface{}{},
		Action:   ai.Action{Type: ai.ActionClarify},
		Response: "Maaf, layanan AI sedang mengalami gangguan. Silakan coba kembali beberapa saat lagi.",
	}
}

// applyRuleBasedOverrides composes DB-backed answers only when the primary AI
// provider is rule-based; Gemini answers are already fact-grounded.
func (h *ChatHandler) applyRuleBasedOverrides(res *ai.AIResult) {
	if !h.ai.IsRuleBased() {
		return
	}
	if h.hotelInfos != nil && isInfoIntent(res.Intent) {
		if infoResp := h.getHotelInfoResponse(res.Intent, res.Language); infoResp != "" {
			res.Response = wrapInfo(infoResp, res.Language)
		} else if res.Language == ai.LangEN {
			res.Response = "Sorry, I don't have that information yet. Please contact Front Office."
		} else {
			res.Response = "Maaf, saya belum memiliki informasi tersebut. Silakan hubungi Front Office."
		}
	}
	if h.recs != nil && isRecommendationIntent(res.Intent) && res.Action.Type == ai.ActionAnswer {
		h.enrichRecommendationResponse(res)
	}
}

// maybeCreateTicket runs the guarded ticket-creation flow per MISSING ROOM NUMBER.md
// and TICKET CONFIRMATION policy. Returns the ticket number when created.
func (h *ChatHandler) maybeCreateTicket(r *http.Request, aiResult *ai.AIResult, originalMessage, resolvedRoom string) *string {
	if !aiResult.RequiresTicket {
		return nil
	}
	// Token room is authoritative: inject into entities so the AI cannot claim another room.
	if resolvedRoom != "" {
		if aiResult.Entities == nil {
			aiResult.Entities = map[string]interface{}{}
		}
		aiResult.Entities["room_number"] = resolvedRoom
	}

	// Missing room flow: ask first, never invent.
	if _, ok := aiResult.Entities["room_number"]; !ok || strings.TrimSpace(fmt.Sprint(aiResult.Entities["room_number"])) == "" {
		return nil
	}

	t, err := h.tickets.CreateFromAI(aiResult, originalMessage)
	if err != nil {
		switch {
		case errors.Is(err, ticketservice.ErrRoomRequired):
			if !strings.Contains(strings.ToLower(aiResult.Response), "kamar") && !strings.Contains(strings.ToLower(aiResult.Response), "room") {
				if aiResult.Language == ai.LangEN {
					aiResult.Response = "Sure, I can help. Could you please tell me your room number?"
				} else {
					aiResult.Response = "Baik, saya bisa membantu. Boleh saya tahu nomor kamar Anda?"
				}
			}
		case ticketservice.IsValidationError(err):
			// Validation error -> no ticket created
		default:
			if aiResult.Language == ai.LangEN {
				aiResult.Response = "Sorry, I couldn't create your ticket due to a temporary issue. Please try again shortly or contact Front Office."
			} else {
				aiResult.Response = "Maaf, saya tidak bisa membuat tiket karena gangguan sementara. Silakan coba lagi atau hubungi Front Office."
			}
		}
		return nil
	}
	s := t.TicketNumber
	return &s
}

// loadHistory returns the last N conversation turns for a session (oldest first),
// mapped to AI Turn roles. Errors are swallowed: history is best-effort context.
func loadHistory(convs *repository.ConversationRepository, sessionID string, n int) []ai.Turn {
	if convs == nil || sessionID == "" {
		return nil
	}
	msgs, err := convs.ListBySession(sessionID, 100)
	if err != nil || len(msgs) == 0 {
		return nil
	}
	// Take only the last n messages
	if len(msgs) > n {
		msgs = msgs[len(msgs)-n:]
	}
	turns := make([]ai.Turn, 0, len(msgs))
	for _, m := range msgs {
		role := "user"
		if m.Role == model.RoleAssistant {
			role = "assistant"
		} else if m.Role != model.RoleUser {
			continue // skip system entries
		}
		turns = append(turns, ai.Turn{Role: role, Content: m.Message})
	}
	return turns
}

// verifiedFacts collects hotel information and local recommendations from the DB
// as plain-text lines the LLM may quote. Never includes internal IDs.
func (h *ChatHandler) verifiedFacts() []string {
	var facts []string
	if h.hotelInfos != nil {
		if items, err := h.hotelInfos.ListActive(nil); err == nil {
			for _, it := range items {
				facts = append(facts, fmt.Sprintf("HOTEL | %s | %s", it.Category, it.Content))
			}
		}
	}
	if h.recs != nil {
		if items, err := h.recs.ListActive(nil, nil); err == nil {
			count := 0
			for _, r := range items {
				if count >= 8 {
					break
				}
				line := fmt.Sprintf("RECOMMENDATION | %s (%s)", r.Name, r.Category)
				if r.Description != nil && *r.Description != "" {
					line += fmt.Sprintf(" - %s", *r.Description)
				}
				if r.Address != nil && *r.Address != "" {
					line += fmt.Sprintf(" | address: %s", *r.Address)
				}
				if r.PriceMin != nil && r.PriceMax != nil && *r.PriceMax > 0 {
					line += fmt.Sprintf(" | price Rp%d-Rp%d", *r.PriceMin, *r.PriceMax)
				}
				if r.DistanceKm != nil {
					line += fmt.Sprintf(" | distance %.1f km", *r.DistanceKm)
				}
				if r.MapsLink != nil && *r.MapsLink != "" {
					line += fmt.Sprintf(" | Google Maps: %s", *r.MapsLink)
				}
				facts = append(facts, line)
				count++
			}
		}
	}
	// Live weather for Pekanbaru (Open-Meteo, no API key, best-effort).
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	if w := currentWeather(ctx); w != "" {
		facts = append(facts, "LIVE_WEATHER | Pekanbaru | "+w)
	}
	return facts
}

// wrapInfo softens template answers from the rule-based provider.
func wrapInfo(fact, lang string) string {
	// zh -> English suffix (the DB fact content itself is Indonesian; the suffix
	// stays neutral so a Chinese speaker reading English is not further confused).
	if lang == ai.LangEN || lang == ai.LangZH {
		return fact + "\n\nIs there anything else I can help you with?"
	}
	return fact + "\n\nAda lagi yang bisa saya bantu?"
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
		if it.MapsLink != nil && *it.MapsLink != "" {
			b.WriteString("\n   📍 ")
			b.WriteString(*it.MapsLink)
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
