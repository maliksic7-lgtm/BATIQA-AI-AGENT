package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"batiqa-ai/internal/model"
	ticketservice "batiqa-ai/internal/service/ticket"
)

// TicketHandler handles /api/tickets per docs/API.md
type TicketHandler struct {
	tickets *ticketservice.Service
}

func NewTicketHandler(svc *ticketservice.Service) *TicketHandler {
	return &TicketHandler{tickets: svc}
}

// Create handles POST /api/tickets per CREATE TICKET.MD
// Request: {room_number, department, category, description, priority}
// Response 201: {ticket_number, room_number, department, category, priority, status}
func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	var req ticketservice.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}
	t, err := h.tickets.Create(req)
	if err != nil {
		// Distinguish validation vs internal
		if ticketservice.IsValidationError(err) || errors.Is(err, ticketservice.ErrRoomRequired) {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create ticket")
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"ticket_number": t.TicketNumber,
		"room_number":   t.RoomNumber,
		"department":    t.Department,
		"category":      t.Category,
		"priority":      t.Priority,
		"status":        t.Status,
	})
}

// List handles GET /api/tickets?department=&status=&priority= per GET TICKETS.MD
func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	q := r.URL.Query()
	filter := model.TicketFilter{Limit: 50}
	if v := strings.TrimSpace(q.Get("department")); v != "" {
		filter.Department = &v
	}
	if v := strings.TrimSpace(q.Get("status")); v != "" {
		filter.Status = &v
	}
	if v := strings.TrimSpace(q.Get("priority")); v != "" {
		filter.Priority = &v
	}
	if v := strings.TrimSpace(q.Get("room_number")); v != "" {
		filter.RoomNumber = &v
	}
	// Also support room query for guest my requests
	if v := strings.TrimSpace(q.Get("room")); v != "" {
		filter.RoomNumber = &v
	}

	list, err := h.tickets.List(filter)
	if err != nil {
		if ticketservice.IsValidationError(err) {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list tickets")
		return
	}
	// Ensure empty slice not null
	if list == nil {
		list = []*model.Ticket{}
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"tickets": list})
}

// GetDetail handles GET /api/tickets/:id per GET TICKETS DETAIL.MD
// :id is ticket_number like TKT-000001
func (h *TicketHandler) GetDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	ticketNumber := extractTicketID(r.URL.Path)
	if ticketNumber == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "ticket_number is required")
		return
	}
	t, err := h.tickets.GetByTicketNumber(ticketNumber)
	if err != nil {
		if errors.Is(err, ticketservice.ErrTicketNotFound) {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if ticketservice.IsValidationError(err) {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get ticket")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ticket_number": t.TicketNumber,
		"room_number":   t.RoomNumber,
		"department":    t.Department,
		"category":      t.Category,
		"description":   t.Description,
		"priority":      t.Priority,
		"status":        t.Status,
		"created_at":    t.CreatedAt,
		"updated_at":    t.UpdatedAt,
		"resolved_at":   t.ResolvedAt,
	})
}

// UpdateStatus handles PATCH /api/tickets/:id/status per UPDATE TICKET STATUS.MD
// Request: {status: "IN_PROGRESS"|"RESOLVED"|"CANCELLED"|"OPEN"}
// Response: {ticket_number, status}
func (h *TicketHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	ticketNumber := extractTicketIDForStatus(r.URL.Path)
	if ticketNumber == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "ticket_number is required")
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}
	if strings.TrimSpace(req.Status) == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "status is required")
		return
	}
	t, err := h.tickets.UpdateStatus(ticketNumber, req.Status)
	if err != nil {
		switch {
		case errors.Is(err, ticketservice.ErrTicketNotFound):
			WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		case errors.Is(err, ticketservice.ErrInvalidTransition) || ticketservice.IsValidationError(err):
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		default:
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update status")
		}
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ticket_number": t.TicketNumber,
		"status":        t.Status,
	})
}

// UpdatePriority handles PATCH /api/tickets/:id/priority (staff only per USER ROLES.md)
// Request: {priority: "LOW"|"MEDIUM"|"HIGH"}
func (h *TicketHandler) UpdatePriority(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	ticketNumber := extractTicketID(r.URL.Path)
	if ticketNumber == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "ticket_number is required")
		return
	}
	var req struct {
		Priority string `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}
	if strings.TrimSpace(req.Priority) == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "priority is required")
		return
	}
	t, err := h.tickets.UpdatePriority(ticketNumber, req.Priority)
	if err != nil {
		switch {
		case errors.Is(err, ticketservice.ErrTicketNotFound):
			WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		case ticketservice.IsValidationError(err):
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		default:
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update priority")
		}
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ticket_number": t.TicketNumber,
		"priority":      t.Priority,
	})
}

func extractTicketID(path string) string {
	// path: /api/tickets/TKT-000001 or /api/tickets/TKT-000001/status
	// For GetDetail: /api/tickets/:id
	trim := strings.TrimPrefix(path, "/api/tickets/")
	trim = strings.TrimSuffix(trim, "/")
	// If contains "/", take first part
	if idx := strings.Index(trim, "/"); idx != -1 {
		trim = trim[:idx]
	}
	return trim
}

func extractTicketIDForStatus(path string) string {
	// /api/tickets/TKT-000001/status -> TKT-000001
	trim := strings.TrimPrefix(path, "/api/tickets/")
	trim = strings.TrimSuffix(trim, "/status")
	return strings.TrimSuffix(trim, "/")
}
