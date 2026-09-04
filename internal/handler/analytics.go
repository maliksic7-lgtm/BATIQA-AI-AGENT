package handler

import (
	"net/http"

	"batiqa-ai/internal/repository"
)

// AnalyticsHandler serves GET /api/analytics (staff only) — operational
// metrics proving business value (tickets/day, resolution speed, hotspots).
type AnalyticsHandler struct {
	tickets *repository.TicketRepository
}

func NewAnalyticsHandler(repo *repository.TicketRepository) *AnalyticsHandler {
	return &AnalyticsHandler{tickets: repo}
}

func (h *AnalyticsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	a, err := h.tickets.GetAnalytics(7)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to compute analytics")
		return
	}
	WriteJSON(w, http.StatusOK, a)
}
