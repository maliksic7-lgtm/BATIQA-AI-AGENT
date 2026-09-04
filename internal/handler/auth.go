package handler

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"batiqa-ai/internal/repository"
)

// StaffAuthHandler handles staff login per USER_ROLES.md & SECURITY PRINCIPLES.md.
// Sessions are persisted in MongoDB (survive restarts, revocable) and login is
// rate-limited per email against brute force.
type StaffAuthHandler struct {
	staffRepo *repository.StaffRepository
	sessions  *repository.StaffSessionRepository
}

type staffSession struct {
	StaffID    int64     `bson:"staff_id"`
	Name       string    `bson:"name"`
	Email      string    `bson:"email"`
	Department string    `bson:"department"`
	ExpiresAt  time.Time `bson:"expires_at"`
}

const sessionTTL = 12 * time.Hour
const tokenTTL = sessionTTL // kept for compatibility

func NewStaffAuthHandler(repo *repository.StaffRepository, sessions *repository.StaffSessionRepository) *StaffAuthHandler {
	return &StaffAuthHandler{staffRepo: repo, sessions: sessions}
}

// ---- Login rate limiting (pluggable) ----
// Default is in-memory (per-instance, unit-testable without a DB). For
// deployments with multiple instances / restarts that must share state, callers
// may install a persistent MongoDB-backed limiter via SetRateLimiter. The plain
// package functions below are kept as the in-memory default used by tests.

const (
	rateWindow   = 15 * time.Minute
	rateMaxFails = 5
)

// rateLimiter abstracts login attempt throttling so the handler can be backed
// by either the in-memory map (default, tests) or a DB store (production).
type rateLimiter interface {
	// limited returns true and a retry duration if the key is over the limit.
	limited(key string) (bool, time.Duration)
	recordFail(key string)
	clear(key string)
}

var activeLimiter rateLimiter = newMemoryLimiter()

type memoryLimiter struct{}

func newMemoryLimiter() *memoryLimiter { return &memoryLimiter{} }

var (
	rateMu     sync.Mutex
	loginFails = map[string][]time.Time{}
)

func (memoryLimiter) limited(key string) (bool, time.Duration) {
	startRateLimiterGC()
	rateMu.Lock()
	defer rateMu.Unlock()
	now := time.Now()
	// Filter the stored fails, keeping only those within the window.
	// NOTE: iterate over loginFails[key] (not the [:0] re-slice, which is empty
	// and would silently wipe every entry on each check).
	fails := loginFails[key][:0]
	for _, t := range loginFails[key] {
		if now.Sub(t) < rateWindow {
			fails = append(fails, t)
		}
	}
	if len(fails) >= rateMaxFails {
		retry := rateWindow - now.Sub(fails[0])
		return true, retry
	}
	loginFails[key] = fails
	return false, 0
}

func (memoryLimiter) recordFail(key string) {
	startRateLimiterGC()
	rateMu.Lock()
	loginFails[key] = append(loginFails[key], time.Now())
	rateMu.Unlock()
}

func (memoryLimiter) clear(key string) {
	rateMu.Lock()
	delete(loginFails, key)
	rateMu.Unlock()
}

// SetRateLimiter installs an alternate (e.g. persistent MongoDB) limiter. A
// nil value restores the in-memory default. Not synchronized: call once at
// startup before serving traffic.
func SetRateLimiter(l rateLimiter) {
	if l == nil {
		activeLimiter = newMemoryLimiter()
	} else {
		activeLimiter = l
	}
}

// mongoRateLimiter backs the persistent path using LoginRateLimitRepository.
// It is only wired when a DB is present (see router.NewWithDB).
type mongoRateLimiter struct {
	repo *repository.LoginRateLimitRepository
}

func NewMongoRateLimiter(repo *repository.LoginRateLimitRepository) *mongoRateLimiter {
	return &mongoRateLimiter{repo: repo}
}

func (m *mongoRateLimiter) limited(key string) (bool, time.Duration) {
	if m.repo == nil {
		return false, 0
	}
	since := time.Now().Add(-rateWindow)
	n, _ := m.repo.CountRecent(key, since)
	if n >= rateMaxFails {
		// Best-effort retry estimate: TTL window remainder (no per-attempt ts
		// stored on the newest slice, so approximate with the full window).
		return true, rateWindow
	}
	return false, 0
}

func (m *mongoRateLimiter) recordFail(key string) {
	if m.repo != nil {
		_ = m.repo.Record(key, time.Now())
	}
}

func (m *mongoRateLimiter) clear(key string) {
	if m.repo != nil {
		m.repo.Clear(key)
	}
}

// startRateLimiterGC periodically purges stale rate-limit entries so the
// in-memory map cannot grow unbounded between requests. Idempotent via sync.Once.
func startRateLimiterGC() {
	rateGCOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(rateWindow)
			defer ticker.Stop()
			for range ticker.C {
				rateMu.Lock()
				cutoff := time.Now().Add(-rateWindow)
				for key, fails := range loginFails {
					kept := fails[:0]
					for _, t := range fails {
						if t.After(cutoff) {
							kept = append(kept, t)
						}
					}
					if len(kept) == 0 {
						delete(loginFails, key)
					} else {
						loginFails[key] = kept
					}
				}
				rateMu.Unlock()
			}
		}()
	})
}

var rateGCOnce sync.Once

func rateLimited(key string) (bool, time.Duration) {
	return activeLimiter.limited(key)
}

func recordFail(key string) {
	activeLimiter.recordFail(key)
}

func clearFails(key string) {
	activeLimiter.clear(key)
}

// ---- DTOs ----

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Staff struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		Department string `json:"department"`
	} `json:"staff"`
}

func (h *StaffAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "email and password required")
		return
	}
	if limited, retry := rateLimited(req.Email); limited {
		w.Header().Set("Retry-After", fmt.Sprintf("%d", int(retry.Seconds())+1))
		WriteError(w, http.StatusTooManyRequests, "RATE_LIMITED", "Terlalu banyak percobaan login. Coba lagi nanti.")
		return
	}

	staff, err := h.staffRepo.FindByEmail(req.Email)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Login failed")
		return
	}
	if staff == nil || bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)) != nil {
		recordFail(req.Email)
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid email or password")
		return
	}
	clearFails(req.Email)

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate token")
		return
	}
	rawToken := hex.EncodeToString(b)
	expires := time.Now().Add(sessionTTL)
	err = h.sessions.Create(r.Context(), repository.StaffSession{
		TokenHash:  hashToken(rawToken),
		StaffID:    staff.ID,
		Name:       staff.Name,
		Email:      staff.Email,
		Department: staff.Department,
		ExpiresAt:  expires,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create session")
		return
	}

	var resp LoginResponse
	resp.Token = rawToken
	resp.Staff.ID = staff.ID
	resp.Staff.Name = staff.Name
	resp.Staff.Email = staff.Email
	resp.Staff.Department = staff.Department

	WriteJSON(w, http.StatusOK, resp)
}

func (h *StaffAuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	sess, ok := GetStaffSession(r)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":         sess.StaffID,
		"name":       sess.Name,
		"email":      sess.Email,
		"department": sess.Department,
	})
}

func (h *StaffAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token != "" {
		_ = h.sessions.Delete(r.Context(), hashToken(token))
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
}

// ---- Token/session helpers ----

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// sessionFromRequest resolves the raw bearer/query token into a persisted session.
func (h *StaffAuthHandler) sessionFromRequest(r *http.Request) (*repository.StaffSession, bool) {
	token := extractToken(r)
	if token == "" {
		token = r.URL.Query().Get("t") // EventSource/img cannot set headers
	}
	if token == "" {
		return nil, false
	}
	sess, err := h.sessions.Find(r.Context(), hashToken(token))
	if err != nil || sess == nil {
		return nil, false
	}
	return sess, true
}

// AuthMiddleware guards staff-only endpoints with a persisted Bearer session.
func (h *StaffAuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, ok := h.sessionFromRequest(r)
		if !ok {
			WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing or invalid token")
			return
		}
		r.Header.Set("X-Staff-ID", fmt.Sprintf("%d", sess.StaffID))
		ctx := context.WithValue(r.Context(), staffCtxKey, staffSession{
			StaffID:    sess.StaffID,
			Name:       sess.Name,
			Email:      sess.Email,
			Department: sess.Department,
			ExpiresAt:  sess.ExpiresAt,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ValidateRequest exposes session validation to other handlers (e.g., SSE).
func (h *StaffAuthHandler) ValidateRequest(r *http.Request) (*repository.StaffSession, bool) {
	return h.sessionFromRequest(r)
}

func extractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(auth)
}

type ctxKey string

const staffCtxKey ctxKey = "staff"

func GetStaffSession(r *http.Request) (staffSession, bool) {
	val := r.Context().Value(staffCtxKey)
	if val == nil {
		return staffSession{}, false
	}
	s, ok := val.(staffSession)
	return s, ok
}
