package model

import "time"

// Staff represents the staff collection
type Staff struct {
	ID           int64     `json:"id" bson:"_id"`
	Name         string    `json:"name" bson:"name"`
	Email        string    `json:"email" bson:"email"`
	PasswordHash string    `json:"-" bson:"password_hash"` // never expose via JSON
	Department   string    `json:"department" bson:"department"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}

// IsValidStaffDepartment validates staff department (includes ADMIN)
func IsValidStaffDepartment(d string) bool {
	return d == DeptHousekeeping || d == DeptEngineering || d == DeptFrontOffice || d == "ADMIN"
}

// TicketAssignment represents the ticket_assignments collection
type TicketAssignment struct {
	ID         int64     `json:"id" bson:"_id"`
	TicketID   int64     `json:"ticket_id" bson:"ticket_id"`
	StaffID    int64     `json:"staff_id" bson:"staff_id"`
	AssignedAt time.Time `json:"assigned_at" bson:"assigned_at"`
}
