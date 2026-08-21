package repository

import (
	"database/sql"
	"fmt"

	"batiqa-ai/internal/model"
)

// AssignmentRepository handles ticket_assignments table.
type AssignmentRepository struct {
	db *sql.DB
}

func NewAssignmentRepository(db *sql.DB) *AssignmentRepository {
	return &AssignmentRepository{db: db}
}

// Assign creates a ticket assignment (staff to ticket). Uses parameterized query.
func (r *AssignmentRepository) Assign(ticketID, staffID int64) (*model.TicketAssignment, error) {
	query := `INSERT INTO ticket_assignments (ticket_id, staff_id) VALUES (?, ?)`
	res, err := r.db.Exec(query, ticketID, staffID)
	if err != nil {
		return nil, fmt.Errorf("assignment Assign: %w", err)
	}
	id, _ := res.LastInsertId()
	// Fetch to get assigned_at
	row := r.db.QueryRow(`SELECT id, ticket_id, staff_id, assigned_at FROM ticket_assignments WHERE id = ?`, id)
	var a model.TicketAssignment
	if err := row.Scan(&a.ID, &a.TicketID, &a.StaffID, &a.AssignedAt); err != nil {
		return nil, fmt.Errorf("assignment scan: %w", err)
	}
	return &a, nil
}

// ListByTicket returns assignments for a ticket.
func (r *AssignmentRepository) ListByTicket(ticketID int64) ([]*model.TicketAssignment, error) {
	rows, err := r.db.Query(`SELECT id, ticket_id, staff_id, assigned_at FROM ticket_assignments WHERE ticket_id = ? ORDER BY assigned_at ASC`, ticketID)
	if err != nil {
		return nil, fmt.Errorf("assignment ListByTicket: %w", err)
	}
	defer rows.Close()

	var out []*model.TicketAssignment
	for rows.Next() {
		var a model.TicketAssignment
		if err := rows.Scan(&a.ID, &a.TicketID, &a.StaffID, &a.AssignedAt); err != nil {
			return nil, fmt.Errorf("assignment scan: %w", err)
		}
		out = append(out, &a)
	}
	return out, rows.Err()
}
