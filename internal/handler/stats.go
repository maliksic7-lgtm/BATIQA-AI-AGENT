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
	TotalOpen     int64 `json:"total_open"`
	HighPriority  int64 `json:"high_priority"`
	Housekeeping  int64 `json:"housekeeping"`
	Engineering   int64 `json:"engineering"`
	ResolvedToday int64 `json:"resolved_today"`
	TotalTickets  int64 `json:"total_tickets"`
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
