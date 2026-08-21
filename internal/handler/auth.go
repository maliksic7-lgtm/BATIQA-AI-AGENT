package handler

import (
	"context"
	"crypto/rand"
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

// StaffAuthHandler handles staff login per USER_ROLES.md & SECURITY PRINCIPLES.md
type StaffAuthHandler struct {
	staffRepo *repository.StaffRepository
}

type staffSession struct {
	StaffID    int64
	Email      string
	Name       string
	Department string
	Expires    time.Time
}

var (
	tokenStore = sync.Map{}
	tokenTTL   = 12 * time.Hour
)

func NewStaffAuthHandler(repo *repository.StaffRepository) *StaffAuthHandler {
	return &StaffAuthHandler{staffRepo: repo}
}

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
	req.Password = strings.TrimSpace(req.Password)
	if req.Email == "" || req.Password == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "email and password required")
		return
	}
	staff, err := h.staffRepo.FindByEmail(req.Email)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Login failed")
		return
	}
	if staff == nil {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid email or password")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)); err != nil {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid email or password")
		return
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate token")
		return
	}
	token := hex.EncodeToString(b)
	session := staffSession{
		StaffID:    staff.ID,
		Email:      staff.Email,
		Name:       staff.Name,
		Department: staff.Department,
		Expires:    time.Now().Add(tokenTTL),
	}
	tokenStore.Store(token, session)

	var resp LoginResponse
	resp.Token = token
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
		tokenStore.Delete(token)
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing token")
			return
		}
		val, ok := tokenStore.Load(token)
		if !ok {
			WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
			return
		}
		sess := val.(staffSession)
		if time.Now().After(sess.Expires) {
			tokenStore.Delete(token)
			WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Token expired")
			return
		}
		r.Header.Set("X-Staff-ID", fmt.Sprintf("%d", sess.StaffID))
		ctx := context.WithValue(r.Context(), staffCtxKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
