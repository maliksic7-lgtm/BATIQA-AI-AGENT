package handler

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// HealthHandler handles GET /api/health per docs/HEALTH CHECK.md
// Response: {"status":"ok"}. When constructed with a DB handle it also reports
// live connectivity ("status":"ok","database":"ok"|"degraded") so Render and
// monitoring tools can observe infrastructure health. Never returns 5xx for a
// down DB: the app still serves static/cached content, so liveness stays 200.
type HealthHandler struct {
	db *mongo.Database
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func NewHealthHandlerWithDB(db *mongo.Database) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	payload := map[string]interface{}{"status": "ok"}
	if h.db != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := h.db.Client().Ping(ctx, nil); err != nil {
			payload["database"] = "degraded"
		} else {
			payload["database"] = "ok"
		}
	}
	WriteOK(w, payload)
}

// HealthCheck is convenience function for router registration (no DB).
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	NewHealthHandler().ServeHTTP(w, r)
}
