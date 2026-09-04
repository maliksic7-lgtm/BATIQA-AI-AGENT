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

// InfographicsHandler serves GET /api/analytics/infographics (staff only) —
// the "most common guest insights" dashboard slices: top complaints, top
// borrowed items, and top asked questions.
type InfographicsHandler struct {
	tickets  *repository.TicketRepository
	convRepo *repository.ConversationRepository
	orders   *repository.RestaurantOrderRepository
}

func NewInfographicsHandler(ticketRepo *repository.TicketRepository, convRepo *repository.ConversationRepository) *InfographicsHandler {
	return &InfographicsHandler{tickets: ticketRepo, convRepo: convRepo}
}

func NewInfographicsHandlerFull(ticketRepo *repository.TicketRepository, convRepo *repository.ConversationRepository, orders *repository.RestaurantOrderRepository) *InfographicsHandler {
	return &InfographicsHandler{tickets: ticketRepo, convRepo: convRepo, orders: orders}
}

func (h *InfographicsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	ig, err := h.tickets.GetInfographics()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to compute infographics")
		return
	}
	if h.convRepo != nil {
		intents, ierr := h.convRepo.TopIntents(5)
		if ierr == nil {
			for _, in := range intents {
				ig.TopAsked = append(ig.TopAsked, repository.CategoryCount{Category: in.Intent, Count: in.Count})
			}
		}
	}
	// Prefer real placed orders for "most ordered"; fall back to chat mentions
	// when the restaurant_orders store is empty or unavailable.
	if h.orders != nil {
		ordered, oerr := h.orders.TopOrderedItems(5)
		if oerr == nil && len(ordered) > 0 {
			for _, c := range ordered {
				ig.TopOrdered = append(ig.TopOrdered, c)
			}
		}
	}
	if len(ig.TopOrdered) == 0 && h.convRepo != nil {
		ordered, oerr := h.convRepo.TopOrderedItems(5)
		if oerr == nil {
			for _, in := range ordered {
				ig.TopOrdered = append(ig.TopOrdered, repository.CategoryCount{Category: in.Intent, Count: in.Count})
			}
		}
	}
	WriteJSON(w, http.StatusOK, ig)
}
