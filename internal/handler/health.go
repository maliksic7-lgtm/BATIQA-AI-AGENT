package handler

import (
	"net/http"
)

// HealthHandler handles GET /api/health per docs/HEALTH CHECK.md
// Response: {"status":"ok"}
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	WriteOK(w, map[string]string{"status": "ok"})
}

// HealthCheck is convenience function for router registration
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	NewHealthHandler().ServeHTTP(w, r)
}
