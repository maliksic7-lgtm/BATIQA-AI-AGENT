package repository

import (
	"database/sql"
	"fmt"

	"batiqa-ai/internal/model"
)

// GuestRepository handles guests table operations.
type GuestRepository struct {
	db *sql.DB
}

func NewGuestRepository(db *sql.DB) *GuestRepository {
	return &GuestRepository{db: db}
}

// FindBySessionID returns guest by session_id (parameterized query).
func (r *GuestRepository) FindBySessionID(sessionID string) (*model.Guest, error) {
	query := `SELECT id, session_id, room_number, language, created_at, updated_at FROM guests WHERE session_id = ?`
	row := r.db.QueryRow(query, sessionID)

	var g model.Guest
	var room sql.NullString
	err := row.Scan(&g.ID, &g.SessionID, &room, &g.Language, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("guest FindBySessionID: %w", err)
	}
	if room.Valid {
		g.RoomNumber = &room.String
	}
	return &g, nil
}

// Upsert creates or updates guest session. Room number never invented; caller must provide validated value.
func (r *GuestRepository) Upsert(sessionID string, roomNumber *string, language string) (*model.Guest, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}
	if language == "" {
		language = "id"
	}
	// Use INSERT ... ON DUPLICATE KEY UPDATE for MySQL
	query := `
		INSERT INTO guests (session_id, room_number, language)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE room_number = VALUES(room_number), language = VALUES(language)
	`
	_, err := r.db.Exec(query, sessionID, roomNumber, language)
	if err != nil {
		return nil, fmt.Errorf("guest Upsert: %w", err)
	}
	return r.FindBySessionID(sessionID)
}

// Create inserts a new guest. Returns error if session_id already exists.
func (r *GuestRepository) Create(g *model.Guest) error {
	query := `INSERT INTO guests (session_id, room_number, language) VALUES (?, ?, ?)`
	res, err := r.db.Exec(query, g.SessionID, g.RoomNumber, g.Language)
	if err != nil {
		return fmt.Errorf("guest Create: %w", err)
	}
	id, _ := res.LastInsertId()
	g.ID = id
	return nil
}
