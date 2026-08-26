package handler

import (
	"net/http"
	"strconv"
	"strings"

	"batiqa-ai/internal/repository"
)

// firstNonEmpty returns the first non-empty string.
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// HotelHandler handles GET /api/hotel-info per HOTEL INFORMATION.MD
type HotelHandler struct {
	repo *repository.HotelInfoRepository
}

func NewHotelHandler(repo *repository.HotelInfoRepository) *HotelHandler {
	return &HotelHandler{repo: repo}
}

func (h *HotelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	q := r.URL.Query()
	var cat *string
	if v := strings.TrimSpace(q.Get("category")); v != "" {
		cat = &v
	}
	items, err := h.repo.ListActive(cat)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch hotel information")
		return
	}
	// Map to response format per HOTEL INFORMATION.MD
	type Item struct {
		ID       int64  `json:"id"`
		Category string `json:"category"`
		Title    string `json:"title"`
		Content  string `json:"content"`
	}
	out := make([]Item, 0, len(items))
	for _, it := range items {
		out = append(out, Item{ID: it.ID, Category: it.Category, Title: it.Title, Content: it.Content})
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"items": out})
}

// RecommendationHandler handles GET /api/recommendations per RECOMMENDATIONS.MD
type RecommendationHandler struct {
	repo *repository.RecommendationRepository
}

func NewRecommendationHandler(repo *repository.RecommendationRepository) *RecommendationHandler {
	return &RecommendationHandler{repo: repo}
}

func (h *RecommendationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	q := r.URL.Query()
	var cat *string
	if v := strings.TrimSpace(q.Get("category")); v != "" {
		cat = &v
	}
	var maxPrice *int
	if raw := firstNonEmpty(strings.TrimSpace(q.Get("max_price")), strings.TrimSpace(q.Get("maxPrice"))); raw != "" {
		// parse int safely (rejects non-numeric, negative, overflow)
		if n, err := strconv.Atoi(raw); err == nil && n >= 0 {
			maxPrice = &n
		}
	}

	items, err := h.repo.ListActive(cat, maxPrice)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch recommendations")
		return
	}
	type Item struct {
		Name       string   `json:"name"`
		Category   string   `json:"category"`
		PriceMin   *int     `json:"price_min"`
		PriceMax   *int     `json:"price_max"`
		DistanceKm *float64 `json:"distance_km"`
		Address    *string  `json:"address"`
	}
	out := make([]Item, 0, len(items))
	for _, it := range items {
		out = append(out, Item{
			Name:       it.Name,
			Category:   it.Category,
			PriceMin:   it.PriceMin,
			PriceMax:   it.PriceMax,
			DistanceKm: it.DistanceKm,
			Address:    it.Address,
		})
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"items": out})
}
