package model

import "time"

// Staff represents staff table
type Staff struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // never expose via JSON
	Department   string    `json:"department" db:"department"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// IsValidStaffDepartment validates staff department (includes ADMIN)
func IsValidStaffDepartment(d string) bool {
	return d == DeptHousekeeping || d == DeptEngineering || d == DeptFrontOffice || d == "ADMIN"
}

// TicketAssignment represents ticket_assignments table
type TicketAssignment struct {
	ID         int64     `json:"id" db:"id"`
	TicketID   int64     `json:"ticket_id" db:"ticket_id"`
	StaffID    int64     `json:"staff_id" db:"staff_id"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
}
