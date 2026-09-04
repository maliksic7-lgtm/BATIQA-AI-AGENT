package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"batiqa-ai/internal/guesttoken"
)

// devBypassRoom is the room auto-assigned to unauthenticated requests when the
// server runs in development mode (ENV=development). This lets the guest app be
// demoed without scanning a QR code, while production still enforces tokens.
const devBypassRoom = "305"

type guestCtxKey string

const guestRoomKey guestCtxKey = "guest_room"

// GuestAuthMiddleware validates a guest token (header X-Guest-Token or ?t= query),
// injecting the verified room into the request context and X-Guest-Room header.
// Requests without a valid token receive 401.
func GuestAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("X-Guest-Token")
		if tok == "" {
			tok = r.URL.Query().Get("t")
		}

		var room string
		if tok != "" {
			parsed, err := guesttoken.Parse(tok, guesttoken.Secret())
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Sesi tidak valid. Silakan scan ulang QR di kamar Anda.")
				return
			}
			room = parsed
		} else if strings.EqualFold(os.Getenv("ENV"), "development") {
			// Dev bypass: allow demos without QR. Towards real deployments this
			// is guarded by ENV, so production still returns 401 for missing tokens.
			room = devBypassRoom
		} else {
			WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Guest token required - please scan the QR code in your room")
			return
		}

		r.Header.Set("X-Guest-Room", room)
		ctx := context.WithValue(r.Context(), guestRoomKey, room)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GuestRoom returns the verified room from context (empty if absent).
func GuestRoom(r *http.Request) string {
	v := r.Context().Value(guestRoomKey)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// EitherAuth allows staff (Bearer) OR guest token; used for shared endpoints
// like ticket listing where staff see everything and guests are scoped later.
func (h *StaffAuthHandler) EitherAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := h.sessionFromRequest(r); ok {
			next.ServeHTTP(w, r)
			return
		}
		GuestAuthMiddleware(next).ServeHTTP(w, r)
	})
}

// QRHandler serves GET /api/rooms/{room}/qr -> PNG QR linking to the guest app
// with a signed token (staff only).
type QRHandler struct {
	auth *StaffAuthHandler
}

func NewQRHandler(auth *StaffAuthHandler) *QRHandler {
	return &QRHandler{auth: auth}
}

func (h *QRHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	room := strings.ToUpper(strings.TrimSpace(extractSegment(r.URL.Path, "/api/rooms/", "/qr")))
	if room == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "room is required")
		return
	}
	token, err := guesttoken.New(room, 90*24*time.Hour, guesttoken.Secret())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate token")
		return
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto == "https" {
		scheme = "https"
	}
	link := fmt.Sprintf("%s://%s/?t=%s", scheme, r.Host, url.QueryEscape(token))

	png, err := qrcodeEncode(link, 512)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to render QR")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(png)
}

// GuestScopedDetail wraps a ticket-detail handler, enforcing that guests can
// only view tickets belonging to their verified room.
func GuestScopedDetail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		guestRoom := GuestRoom(r)
		if guestRoom == "" {
			WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Guest token required")
			return
		}
		// Peek at the ticket room via a recording wrapper: simplest correct way
		// is to let the inner handler run but verify first through the service.
		// We re-implement the guard by extracting the ticket number and checking ownership.
		ticketNumber := extractSegment(r.URL.Path, "/api/tickets/", "")
		if ticketNumber == "" {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "ticket_number is required")
			return
		}
		if !guestTicketOwnerOK(r, ticketNumber, guestRoom) {
			WriteError(w, http.StatusForbidden, "FORBIDDEN", "This ticket does not belong to your room")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// guestTicketOwnerOK checks ticket ownership using the shared ticket lookup hook.
var guestTicketOwnerOK = defaultGuestTicketOwner

// SetTicketOwnershipChecker lets the router inject a DB-backed owner check.
func SetTicketOwnershipChecker(fn func(r *http.Request, ticketNumber, room string) bool) {
	guestTicketOwnerOK = fn
}

func defaultGuestTicketOwner(r *http.Request, ticketNumber, room string) bool { return false }

// extractSegment pulls the middle segment between prefix and suffix.
func extractSegment(path, prefix, suffix string) string {
	s := strings.TrimPrefix(path, prefix)
	return strings.TrimSuffix(s, suffix)
}
