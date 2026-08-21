package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"batiqa-ai/internal/model"
	"batiqa-ai/internal/repository"
	"batiqa-ai/internal/service/ai"
	ticketservice "batiqa-ai/internal/service/ticket"
)

// ChatHandler handles POST /api/chat per AI CHAT.MD and USER_FLOW.md
// Flow: Guest Message -> AI -> Structured AI Result -> Backend Validation -> Ticket Service -> Repository -> DB
// For info intents, it retrieves verified hotel information from DB per AI KNOWLEDGE SOURCE.md
type ChatHandler struct {
	ai         *ai.Service
	tickets    *ticketservice.Service
	convs      *repository.ConversationRepository
	guests     *repository.GuestRepository
	hotelInfos *repository.HotelInfoRepository
}

func NewChatHandler(aiSvc *ai.Service, ticketSvc *ticketservice.Service, convRepo *repository.ConversationRepository, guestRepo *repository.GuestRepository) *ChatHandler {
	return &ChatHandler{ai: aiSvc, tickets: ticketSvc, convs: convRepo, guests: guestRepo}
}

func NewChatHandlerWithHotel(aiSvc *ai.Service, ticketSvc *ticketservice.Service, convRepo *repository.ConversationRepository, guestRepo *repository.GuestRepository, hotelRepo *repository.HotelInfoRepository) *ChatHandler {
	return &ChatHandler{ai: aiSvc, tickets: ticketSvc, convs: convRepo, guests: guestRepo, hotelInfos: hotelRepo}
}

// ChatRequest per AI CHAT.MD
type ChatRequest struct {
	SessionID  string `json:"session_id"`
	RoomNumber string `json:"room_number"`
	Message    string `json:"message"`
}

// ChatResponse per AI CHAT.MD
type ChatResponse struct {
	Message        string `json:"message"`
	Intent         string `json:"intent"`
	RequiresTicket bool   `json:"requires_ticket"`
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

	// Persist guest session (upsert) - room not invented
	var roomPtr *string
	if req.RoomNumber != "" {
		roomPtr = &req.RoomNumber
	}
	// Language will be detected by AI
	lang := ai.DetectLanguage(req.Message)
	if _, err := h.guests.Upsert(req.SessionID, roomPtr, lang); err != nil {
		// Non-fatal, log but continue (guest table is optional)
	}

	// Store user conversation
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
	if err != nil {
		if aiResult == nil {
			aiResult = &ai.AIResult{
				Intent:   ai.IntentUnknown,
				Language: ai.LangID,
				Entities: map[string]interface{}{},
				Action:   ai.Action{Type: ai.ActionClarify},
				Response: "Maaf, layanan AI sedang mengalami gangguan. Silakan coba kembali beberapa saat lagi.",
			}
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
			// Do not create ticket, keep requires_ticket true but ticket_id null, response already asks for room
			// Ensure response is ask for room
			requiresTicket = true
			ticketID = nil
		} else {
				t, err := h.tickets.CreateFromAI(aiResult, req.Message)
			if err != nil {
				// Distinguish validation vs internal DB error
				lowerErr := strings.ToLower(err.Error())
				if strings.Contains(lowerErr, "room_number is required") {
					requiresTicket = true
					ticketID = nil
					if !strings.Contains(strings.ToLower(aiResult.Response), "kamar") && !strings.Contains(strings.ToLower(aiResult.Response), "room") {
						if aiResult.Language == ai.LangEN {
							aiResult.Response = "Sure, I can help. Could you please tell me your room number?"
						} else {
							aiResult.Response = "Baik, saya bisa membantu. Boleh saya tahu nomor kamar Anda?"
						}
					}
				} else if strings.Contains(lowerErr, "invalid") || strings.Contains(lowerErr, "required") || strings.Contains(lowerErr, "must") {
					// Validation error -> do not create ticket
					requiresTicket = false
					ticketID = nil
				} else {
					// Internal DB error (e.g., ticket create: ...) -> keep requires_ticket true but no ticket, inform guest of temporary failure
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

	// Also persist ticket creation conversation if needed? Already done

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



