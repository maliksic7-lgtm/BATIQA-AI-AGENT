package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"batiqa-ai/internal/events"
)

// SSEHandler streams ticket lifecycle events to the staff dashboard via
// Server-Sent Events (GET /api/events).
type SSEHandler struct {
	auth   *StaffAuthHandler
	broker *events.Broker
}

func NewSSEHandler(auth *StaffAuthHandler, broker *events.Broker) *SSEHandler {
	return &SSEHandler{auth: auth, broker: broker}
}

func (h *SSEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	if _, ok := h.auth.ValidateRequest(r); !ok {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Staff authentication required")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Streaming not supported")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	subID, ch := h.broker.Subscribe()
	defer h.broker.Unsubscribe(subID)

	// Initial hello so the client knows the stream is live
	fmt.Fprint(w, "event: ready\ndata: {\"status\":\"live\"}\n\n")
	flusher.Flush()

	ping := time.NewTicker(25 * time.Second)
	defer ping.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ping.C:
			fmt.Fprint(w, ": ping\n\n")
			flusher.Flush()
		case ev, open := <-ch:
			if !open {
				return
			}
			data, err := json.Marshal(ev.Data)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "event: %s\ndata: {\"type\":%q,\"ticket\":%s}\n\n", ev.Type, ev.Type, data)
			flusher.Flush()
		}
	}
}

var _ = context.Background
