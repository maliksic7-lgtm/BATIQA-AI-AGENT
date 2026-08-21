package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"batiqa-ai/internal/model"
)

// TicketRepository handles tickets table with parameterized queries.
type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

// Create inserts a ticket and generates ticket_number as TKT-{id} (zero-padded 6 digits).
// It validates enums before insert per DATABASE.md rules.
func (r *TicketRepository) Create(t *model.Ticket) error {
	if !model.IsValidDepartment(t.Department) {
		return fmt.Errorf("invalid department: %s", t.Department)
	}
	if !model.IsValidPriority(t.Priority) {
		return fmt.Errorf("invalid priority: %s", t.Priority)
	}
	if !model.IsValidStatus(t.Status) {
		return fmt.Errorf("invalid status: %s", t.Status)
	}
	if strings.TrimSpace(t.RoomNumber) == "" {
		return fmt.Errorf("room_number is required")
	}
	if strings.TrimSpace(t.Category) == "" {
		return fmt.Errorf("category is required")
	}
	if strings.TrimSpace(t.Description) == "" {
		return fmt.Errorf("description is required")
	}
	// Default status/priority if empty
	if t.Status == "" {
		t.Status = model.StatusOpen
	}
	if t.Priority == "" {
		t.Priority = model.PriorityMedium
	}

	// Use transaction to ensure ticket_number generation is atomic.
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("ticket Create begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Insert with placeholder ticket_number (will be updated)
	placeholder := "TMP"
	query := `INSERT INTO tickets (ticket_number, room_number, department, category, description, priority, status)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := tx.Exec(query, placeholder, t.RoomNumber, t.Department, t.Category, t.Description, t.Priority, t.Status)
	if err != nil {
		return fmt.Errorf("ticket Create insert: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("ticket Create last insert id: %w", err)
	}
	t.ID = id
	ticketNumber := fmt.Sprintf("TKT-%06d", id)

	// Update ticket_number
	if _, err := tx.Exec(`UPDATE tickets SET ticket_number = ? WHERE id = ?`, ticketNumber, id); err != nil {
		return fmt.Errorf("ticket Create update number: %w", err)
	}
	t.TicketNumber = ticketNumber

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ticket Create commit: %w", err)
	}

	// Refresh timestamps
	return r.reload(t)
}

func (r *TicketRepository) reload(t *model.Ticket) error {
	query := `SELECT id, ticket_number, room_number, department, category, description, priority, status, created_at, updated_at, resolved_at
	          FROM tickets WHERE id = ?`
	row := r.db.QueryRow(query, t.ID)
	var resolved sql.NullTime
	err := row.Scan(&t.ID, &t.TicketNumber, &t.RoomNumber, &t.Department, &t.Category, &t.Description, &t.Priority, &t.Status, &t.CreatedAt, &t.UpdatedAt, &resolved)
	if err != nil {
		return fmt.Errorf("ticket reload: %w", err)
	}
	if resolved.Valid {
		t.ResolvedAt = &resolved.Time
	}
	return nil
}

// FindByTicketNumber returns ticket by ticket_number.
func (r *TicketRepository) FindByTicketNumber(ticketNumber string) (*model.Ticket, error) {
	query := `SELECT id, ticket_number, room_number, department, category, description, priority, status, created_at, updated_at, resolved_at
	          FROM tickets WHERE ticket_number = ?`
	row := r.db.QueryRow(query, ticketNumber)
	var t model.Ticket
	var resolved sql.NullTime
	err := row.Scan(&t.ID, &t.TicketNumber, &t.RoomNumber, &t.Department, &t.Category, &t.Description, &t.Priority, &t.Status, &t.CreatedAt, &t.UpdatedAt, &resolved)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("ticket FindByTicketNumber: %w", err)
	}
	if resolved.Valid {
		t.ResolvedAt = &resolved.Time
	}
	return &t, nil
}

// FindByID returns ticket by id.
func (r *TicketRepository) FindByID(id int64) (*model.Ticket, error) {
	query := `SELECT id, ticket_number, room_number, department, category, description, priority, status, created_at, updated_at, resolved_at
	          FROM tickets WHERE id = ?`
	row := r.db.QueryRow(query, id)
	var t model.Ticket
	var resolved sql.NullTime
	err := row.Scan(&t.ID, &t.TicketNumber, &t.RoomNumber, &t.Department, &t.Category, &t.Description, &t.Priority, &t.Status, &t.CreatedAt, &t.UpdatedAt, &resolved)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("ticket FindByID: %w", err)
	}
	if resolved.Valid {
		t.ResolvedAt = &resolved.Time
	}
	return &t, nil
}

// List returns tickets with optional filters, ordered by created_at DESC.
func (r *TicketRepository) List(f model.TicketFilter) ([]*model.Ticket, error) {
	// Build query with parameterized filters (no string concatenation of values)
	base := `SELECT id, ticket_number, room_number, department, category, description, priority, status, created_at, updated_at, resolved_at FROM tickets WHERE 1=1`
	args := []interface{}{}

	if f.Department != nil && *f.Department != "" {
		if !model.IsValidDepartment(*f.Department) {
			return nil, fmt.Errorf("invalid department filter: %s", *f.Department)
		}
		base += ` AND department = ?`
		args = append(args, *f.Department)
	}
	if f.Status != nil && *f.Status != "" {
		if !model.IsValidStatus(*f.Status) {
			return nil, fmt.Errorf("invalid status filter: %s", *f.Status)
		}
		base += ` AND status = ?`
		args = append(args, *f.Status)
	}
	if f.Priority != nil && *f.Priority != "" {
		if !model.IsValidPriority(*f.Priority) {
			return nil, fmt.Errorf("invalid priority filter: %s", *f.Priority)
		}
		base += ` AND priority = ?`
		args = append(args, *f.Priority)
	}
	if f.RoomNumber != nil && *f.RoomNumber != "" {
		base += ` AND room_number = ?`
		args = append(args, *f.RoomNumber)
	}
	base += ` ORDER BY created_at DESC`

	if f.Limit > 0 {
		if f.Limit > 100 {
			f.Limit = 100
		}
		base += ` LIMIT ?`
		args = append(args, f.Limit)
		if f.Offset > 0 {
			base += ` OFFSET ?`
			args = append(args, f.Offset)
		}
	}

	rows, err := r.db.Query(base, args...)
	if err != nil {
		return nil, fmt.Errorf("ticket List: %w", err)
	}
	defer rows.Close()

	var out []*model.Ticket
	for rows.Next() {
		var t model.Ticket
		var resolved sql.NullTime
		if err := rows.Scan(&t.ID, &t.TicketNumber, &t.RoomNumber, &t.Department, &t.Category, &t.Description, &t.Priority, &t.Status, &t.CreatedAt, &t.UpdatedAt, &resolved); err != nil {
			return nil, fmt.Errorf("ticket List scan: %w", err)
		}
		if resolved.Valid {
			t.ResolvedAt = &resolved.Time
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}

// UpdateStatus updates ticket status with validation and sets resolved_at when RESOLVED.
func (r *TicketRepository) UpdateStatus(ticketNumber, newStatus string) (*model.Ticket, error) {
	if !model.IsValidStatus(newStatus) {
		return nil, fmt.Errorf("invalid status: %s", newStatus)
	}
	t, err := r.FindByTicketNumber(ticketNumber)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, fmt.Errorf("ticket not found: %s", ticketNumber)
	}
	if !model.IsValidStatusTransition(t.Status, newStatus) {
		return nil, fmt.Errorf("invalid status transition %s -> %s", t.Status, newStatus)
	}

	var query string
	var args []interface{}
	if newStatus == model.StatusResolved {
		query = `UPDATE tickets SET status = ?, resolved_at = CURRENT_TIMESTAMP WHERE ticket_number = ?`
		args = []interface{}{newStatus, ticketNumber}
	} else {
		// Clear resolved_at if moving away from RESOLVED
		query = `UPDATE tickets SET status = ?, resolved_at = NULL WHERE ticket_number = ?`
		args = []interface{}{newStatus, ticketNumber}
		if newStatus == model.StatusCancelled {
			// keep resolved_at NULL
		}
	}
	if _, err := r.db.Exec(query, args...); err != nil {
		return nil, fmt.Errorf("ticket UpdateStatus: %w", err)
	}
	return r.FindByTicketNumber(ticketNumber)
}

// UpdatePriority updates ticket priority.
func (r *TicketRepository) UpdatePriority(ticketNumber, newPriority string) (*model.Ticket, error) {
	if !model.IsValidPriority(newPriority) {
		return nil, fmt.Errorf("invalid priority: %s", newPriority)
	}
	if _, err := r.db.Exec(`UPDATE tickets SET priority = ? WHERE ticket_number = ?`, newPriority, ticketNumber); err != nil {
		return nil, fmt.Errorf("ticket UpdatePriority: %w", err)
	}
	return r.FindByTicketNumber(ticketNumber)
}

// TicketStats holds dashboard stats
type TicketStats struct {
	Total         int
	Open          int
	High          int
	Housekeeping  int
	Engineering   int
	ResolvedToday int
}

// GetStats returns dashboard statistics for staff overview
func (r *TicketRepository) GetStats() (*TicketStats, error) {
	stats := &TicketStats{}
	queries := map[string]*int{
		`SELECT COUNT(*) FROM tickets`: &stats.Total,
		`SELECT COUNT(*) FROM tickets WHERE status = 'OPEN'`: &stats.Open,
		`SELECT COUNT(*) FROM tickets WHERE priority = 'HIGH' AND status != 'RESOLVED' AND status != 'CANCELLED'`: &stats.High,
		`SELECT COUNT(*) FROM tickets WHERE department = 'HOUSEKEEPING' AND status != 'RESOLVED' AND status != 'CANCELLED'`: &stats.Housekeeping,
		`SELECT COUNT(*) FROM tickets WHERE department = 'ENGINEERING' AND status != 'RESOLVED' AND status != 'CANCELLED'`: &stats.Engineering,
		`SELECT COUNT(*) FROM tickets WHERE status = 'RESOLVED' AND DATE(resolved_at) = CURDATE()`: &stats.ResolvedToday,
	}
	for q, ptr := range queries {
		if err := r.db.QueryRow(q).Scan(ptr); err != nil {
			return nil, fmt.Errorf("stats query %q: %w", q, err)
		}
	}
	return stats, nil
}
