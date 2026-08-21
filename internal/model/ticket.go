package model

import "time"

// Department enum per TICKET ROUTING.md & DATABASE.md
const (
	DeptHousekeeping = "HOUSEKEEPING"
	DeptEngineering  = "ENGINEERING"
	DeptFrontOffice  = "FRONT_OFFICE"
)

// Priority enum per PRIORITY CLASSIFICATION.md
const (
	PriorityLow    = "LOW"
	PriorityMedium = "MEDIUM"
	PriorityHigh   = "HIGH"
)

// Status enum per TICKET LIFECYCLE.md
const (
	StatusOpen       = "OPEN"
	StatusInProgress = "IN_PROGRESS"
	StatusResolved   = "RESOLVED"
	StatusCancelled  = "CANCELLED"
)

// Ticket represents tickets table (main operational table)
type Ticket struct {
	ID           int64      `json:"id" db:"id"`
	TicketNumber string     `json:"ticket_number" db:"ticket_number"`
	RoomNumber   string     `json:"room_number" db:"room_number"`
	Department   string     `json:"department" db:"department"`
	Category     string     `json:"category" db:"category"`
	Description  string     `json:"description" db:"description"`
	Priority     string     `json:"priority" db:"priority"`
	Status       string     `json:"status" db:"status"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	ResolvedAt   *time.Time `json:"resolved_at,omitempty" db:"resolved_at"`
}

// IsValidDepartment validates department
func IsValidDepartment(d string) bool {
	return d == DeptHousekeeping || d == DeptEngineering || d == DeptFrontOffice
}

// IsValidPriority validates priority
func IsValidPriority(p string) bool {
	return p == PriorityLow || p == PriorityMedium || p == PriorityHigh
}

// IsValidStatus validates status
func IsValidStatus(s string) bool {
	return s == StatusOpen || s == StatusInProgress || s == StatusResolved || s == StatusCancelled
}

// IsValidStatusTransition validates allowed transitions (simple: any except RESOLVED->OPEN without CANCELLED)
func IsValidStatusTransition(from, to string) bool {
	if !IsValidStatus(from) || !IsValidStatus(to) {
		return false
	}
	// Allow all transitions in Phase 2; stricter rules can be added later
	return true
}

// TicketFilter for GET /api/tickets?department=&status=&priority=
type TicketFilter struct {
	Department *string
	Status     *string
	Priority   *string
	RoomNumber *string
	Limit      int
	Offset     int
}
