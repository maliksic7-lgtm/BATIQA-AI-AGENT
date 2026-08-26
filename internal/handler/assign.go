package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"batiqa-ai/internal/repository"
	ticketservice "batiqa-ai/internal/service/ticket"
)

// AssignHandler handles ticket assignment to staff (USER ROLES.md).
// POST /api/tickets/{ticket_number}/assign  {staff_id} -> assign (staff only)
// GET  /api/tickets/{ticket_number}/assignments         -> list assignments
type AssignHandler struct {
	tickets *ticketservice.Service
	staff   *repository.StaffRepository
	assigns *repository.AssignmentRepository
}

func NewAssignHandler(ticketSvc *ticketservice.Service, staffRepo *repository.StaffRepository, assignRepo *repository.AssignmentRepository) *AssignHandler {
	return &AssignHandler{tickets: ticketSvc, staff: staffRepo, assigns: assignRepo}
}

func (h *AssignHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleAssign(w, r)
	case http.MethodGet:
		h.list(w, r)
	default:
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
	}
}

func (h *AssignHandler) handleAssign(w http.ResponseWriter, r *http.Request) {
	ticketNumber := extractTicketID(r.URL.Path)
	if ticketNumber == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "ticket_number is required")
		return
	}
	var req struct {
		StaffID int64 `json:"staff_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}
	req.StaffID = resolveStaffID(req.StaffID, r.Header.Get("X-Staff-ID"))
	if req.StaffID <= 0 {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "staff_id is required")
		return
	}
	t, err := h.tickets.GetByTicketNumber(ticketNumber)
	if err != nil {
		writeTicketError(w, err)
		return
	}
	staff, err := h.staff.FindByID(req.StaffID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to find staff")
		return
	}
	if staff == nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "staff not found")
		return
	}
	a, err := h.assigns.Assign(t.ID, req.StaffID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to assign ticket")
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"ticket_number": t.TicketNumber,
		"staff_id":      a.StaffID,
		"assigned_at":   a.AssignedAt,
	})
}

func (h *AssignHandler) list(w http.ResponseWriter, r *http.Request) {
	ticketNumber := extractTicketID(r.URL.Path)
	if ticketNumber == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "ticket_number is required")
		return
	}
	t, err := h.tickets.GetByTicketNumber(ticketNumber)
	if err != nil {
		writeTicketError(w, err)
		return
	}
	list, err := h.assigns.ListByTicket(t.ID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch assignments")
		return
	}
	if list == nil {
		list = nil
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"assignments": list})
}

func writeTicketError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ticketservice.ErrTicketNotFound):
		WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
	case ticketservice.IsValidationError(err):
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
	default:
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get ticket")
	}
}

// resolveStaffID prefers an explicit staff_id from the request body
// (admin assigning on behalf of others) and falls back to the
// authenticated staff from the auth middleware.
func resolveStaffID(bodyID int64, headerVal string) int64 {
	if bodyID > 0 {
		return bodyID
	}
	if v := strings.TrimSpace(headerVal); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			return n
		}
	}
	return 0
}
