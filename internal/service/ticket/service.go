package ticket

import (
	"fmt"
	"strings"

	"batiqa-ai/internal/model"
	"batiqa-ai/internal/repository"
	"batiqa-ai/internal/service/ai"
)

// Service handles ticket business logic per docs/USER_FLOW.md & TICKET LIFECYCLE.md
// Flow: AI Result -> Backend Validation -> Ticket Service -> Repository -> DB
// AI never directly touches DB.
type Service struct {
	tickets *repository.TicketRepository
	guests  *repository.GuestRepository
}

func NewService(ticketRepo *repository.TicketRepository, guestRepo *repository.GuestRepository) *Service {
	return &Service{tickets: ticketRepo, guests: guestRepo}
}

// CreateRequest per docs/CREATE TICKET.MD
type CreateRequest struct {
	RoomNumber  string `json:"room_number"`
	Department  string `json:"department"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// ValidateCreate ensures minimum data per REQUIRED TICKET DATA.md
func ValidateCreate(req CreateRequest) error {
	if strings.TrimSpace(req.RoomNumber) == "" {
		return fmt.Errorf("room_number is required")
	}
	if !isValidRoomNumber(req.RoomNumber) {
		return fmt.Errorf("invalid room_number format: %s", req.RoomNumber)
	}
	if strings.TrimSpace(req.Department) == "" {
		return fmt.Errorf("department is required")
	}
	if !model.IsValidDepartment(strings.ToUpper(req.Department)) {
		return fmt.Errorf("invalid department: %s", req.Department)
	}
	if strings.TrimSpace(req.Category) == "" {
		return fmt.Errorf("category is required")
	}
	if !ai.IsValidIntent(strings.ToUpper(req.Category)) {
		return fmt.Errorf("invalid category (intent): %s", req.Category)
	}
	if strings.TrimSpace(req.Description) == "" {
		return fmt.Errorf("description is required")
	}
	if strings.TrimSpace(req.Priority) == "" {
		return fmt.Errorf("priority is required")
	}
	if !model.IsValidPriority(strings.ToUpper(req.Priority)) {
		return fmt.Errorf("invalid priority: %s", req.Priority)
	}
	return nil
}

func isValidRoomNumber(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// Simple: 2-4 digits optional letter, same as AI validator
	if len(s) < 2 || len(s) > 5 {
		return false
	}
	for i, r := range s {
		if i < len(s)-1 || (r >= '0' && r <= '9') {
			if r < '0' || r > '9' {
				if i != len(s)-1 {
					return false
				}
				if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
					return false
				}
			}
		} else {
			// last char may be letter
			if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') {
				return false
			}
		}
	}
	// must start with digit
	if s[0] < '0' || s[0] > '9' {
		return false
	}
	return true
}

// Create validates and creates ticket. Returns created ticket or error.
func (s *Service) Create(req CreateRequest) (*model.Ticket, error) {
	// Normalize
	req.Department = strings.ToUpper(strings.TrimSpace(req.Department))
	req.Category = strings.ToUpper(strings.TrimSpace(req.Category))
	req.Priority = strings.ToUpper(strings.TrimSpace(req.Priority))
	req.RoomNumber = strings.TrimSpace(req.RoomNumber)

	if err := ValidateCreate(req); err != nil {
		return nil, err
	}

	t := &model.Ticket{
		RoomNumber:  req.RoomNumber,
		Department:  req.Department,
		Category:    req.Category,
		Description: strings.TrimSpace(req.Description),
		Priority:    req.Priority,
		Status:      model.StatusOpen,
	}
	if err := s.tickets.Create(t); err != nil {
		return nil, fmt.Errorf("ticket create: %w", err)
	}
	return t, nil
}

// CreateFromAI creates ticket from validated AIResult.
// It ensures AI cannot directly create DB without backend validation.
// Requires aiResult validated via ai.Validate, and checks minimum data.
func (s *Service) CreateFromAI(aiResult *ai.AIResult, originalMessage string) (*model.Ticket, error) {
	if aiResult == nil {
		return nil, fmt.Errorf("nil AI result")
	}
	if !aiResult.RequiresTicket {
		return nil, fmt.Errorf("intent %s does not require ticket", aiResult.Intent)
	}
	if aiResult.Action.Type != ai.ActionCreateTicket {
		return nil, fmt.Errorf("action type %s not CREATE_TICKET", aiResult.Action.Type)
	}

	// Extract required fields
	room := ""
	if v, ok := aiResult.Entities["room_number"]; ok {
		if s, ok := v.(string); ok {
			room = strings.TrimSpace(s)
		}
	}
	if room == "" {
		return nil, fmt.Errorf("room_number is required")
	}
	dept := strings.ToUpper(aiResult.Action.Department)
	priority := strings.ToUpper(aiResult.Action.Priority)
	category := aiResult.Intent
	description := buildDescription(aiResult, originalMessage)

	req := CreateRequest{
		RoomNumber:  room,
		Department:  dept,
		Category:    category,
		Description: description,
		Priority:    priority,
	}
	return s.Create(req)
}

func buildDescription(aiResult *ai.AIResult, originalMessage string) string {
	// Prefer problem entity, then item+quantity, then original message
	if v, ok := aiResult.Entities["problem"]; ok {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	if item, ok := aiResult.Entities["item"]; ok {
		if s, ok := item.(string); ok {
			qty := ""
			if q, ok := aiResult.Entities["quantity"]; ok {
				switch val := q.(type) {
				case int:
					qty = fmt.Sprintf("%d ", val)
				case float64:
					qty = fmt.Sprintf("%d ", int(val))
				}
			}
			return fmt.Sprintf("%s%s", qty, s)
		}
	}
	// Fallback to original message trimmed
	msg := strings.TrimSpace(originalMessage)
	if len(msg) > 500 {
		msg = msg[:500]
	}
	if msg == "" {
		msg = aiResult.Intent
	}
	return msg
}

// List returns tickets with filtering per GET TICKETS.MD
func (s *Service) List(filter model.TicketFilter) ([]*model.Ticket, error) {
	// Validate filter enums if set
	if filter.Department != nil && *filter.Department != "" {
		dept := strings.ToUpper(*filter.Department)
		if !model.IsValidDepartment(dept) {
			return nil, fmt.Errorf("invalid department filter: %s", *filter.Department)
		}
		filter.Department = &dept
	}
	if filter.Status != nil && *filter.Status != "" {
		st := strings.ToUpper(*filter.Status)
		if !model.IsValidStatus(st) {
			return nil, fmt.Errorf("invalid status filter: %s", *filter.Status)
		}
		filter.Status = &st
	}
	if filter.Priority != nil && *filter.Priority != "" {
		p := strings.ToUpper(*filter.Priority)
		if !model.IsValidPriority(p) {
			return nil, fmt.Errorf("invalid priority filter: %s", *filter.Priority)
		}
		filter.Priority = &p
	}
	return s.tickets.List(filter)
}

// GetByTicketNumber returns ticket detail per GET TICKETS DETAIL.MD
func (s *Service) GetByTicketNumber(ticketNumber string) (*model.Ticket, error) {
	if strings.TrimSpace(ticketNumber) == "" {
		return nil, fmt.Errorf("ticket_number is required")
	}
	t, err := s.tickets.FindByTicketNumber(strings.TrimSpace(ticketNumber))
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, fmt.Errorf("ticket not found: %s", ticketNumber)
	}
	return t, nil
}

// UpdateStatus validates transition OPEN->IN_PROGRESS->RESOLVED, optional CANCELLED per TICKET LIFECYCLE.md
func (s *Service) UpdateStatus(ticketNumber, newStatus string) (*model.Ticket, error) {
	newStatus = strings.ToUpper(strings.TrimSpace(newStatus))
	if !model.IsValidStatus(newStatus) {
		return nil, fmt.Errorf("invalid status: %s", newStatus)
	}
	// Ensure ticket exists and transition is valid
	current, err := s.GetByTicketNumber(ticketNumber)
	if err != nil {
		return nil, err
	}
	// Validate transition: RESOLVED/CANCELLED are terminal (no outgoing), per TICKET LIFECYCLE.md
	if current.Status == model.StatusResolved || current.Status == model.StatusCancelled {
		if newStatus != current.Status {
			return nil, fmt.Errorf("cannot transition from %s to %s", current.Status, newStatus)
		}
		return current, nil
	}
	if current.Status == model.StatusOpen && newStatus == model.StatusOpen {
		return current, nil
	}
	return s.tickets.UpdateStatus(ticketNumber, newStatus)
}

// UpdatePriority validates and updates priority
func (s *Service) UpdatePriority(ticketNumber, newPriority string) (*model.Ticket, error) {
	newPriority = strings.ToUpper(strings.TrimSpace(newPriority))
	if !model.IsValidPriority(newPriority) {
		return nil, fmt.Errorf("invalid priority: %s", newPriority)
	}
	_, err := s.GetByTicketNumber(ticketNumber)
	if err != nil {
		return nil, err
	}
	return s.tickets.UpdatePriority(ticketNumber, newPriority)
}
