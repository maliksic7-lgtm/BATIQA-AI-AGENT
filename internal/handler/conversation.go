package handler

import (
	"net/http"
	"strconv"
	"strings"

	"batiqa-ai/internal/repository"
)

// ConversationHandler handles GET /api/conversations?session_id=&limit=
// Returns chat history for a guest session (oldest first, capped 100).
type ConversationHandler struct {
	convs *repository.ConversationRepository
}

func NewConversationHandler(repo *repository.ConversationRepository) *ConversationHandler {
	return &ConversationHandler{convs: repo}
}

func (h *ConversationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	sessionID := strings.TrimSpace(r.URL.Query().Get("session_id"))
	if sessionID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "session_id is required")
		return
	}
	limit := 50
	if v := strings.TrimSpace(r.URL.Query().Get("limit")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	msgs, err := h.convs.ListBySession(sessionID, limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch conversations")
		return
	}
	type Message struct {
		ID        int64   `json:"id"`
		Role      string  `json:"role"`
		Message   string  `json:"message"`
		Intent    *string `json:"intent,omitempty"`
		CreatedAt string  `json:"created_at"`
	}
	out := make([]Message, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, Message{
			ID:        m.ID,
			Role:      m.Role,
			Message:   m.Message,
			Intent:    m.Intent,
			CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"messages": out})
}
