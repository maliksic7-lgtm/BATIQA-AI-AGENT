package handler

import (
	"net/http"

	"batiqa-ai/internal/repository"
)

// StatsHandler handles GET /api/tickets/stats for staff dashboard overview
type StatsHandler struct {
	tickets *repository.TicketRepository
}

func NewStatsHandler(repo *repository.TicketRepository) *StatsHandler {
	return &StatsHandler{tickets: repo}
}

type StatsResponse struct {
	TotalOpen     int `json:"total_open"`
	HighPriority  int `json:"high_priority"`
	Housekeeping  int `json:"housekeeping"`
	Engineering   int `json:"engineering"`
	ResolvedToday int `json:"resolved_today"`
	TotalTickets  int `json:"total_tickets"`
}

func (h *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	stats, err := h.tickets.GetStats()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch stats")
		return
	}
	resp := StatsResponse{
		TotalOpen:     stats.Open,
		HighPriority:  stats.High,
		Housekeeping:  stats.Housekeeping,
		Engineering:   stats.Engineering,
		ResolvedToday: stats.ResolvedToday,
		TotalTickets:  stats.Total,
	}
	WriteJSON(w, http.StatusOK, resp)
}
