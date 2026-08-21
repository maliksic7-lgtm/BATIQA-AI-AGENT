package repository

import (
	"database/sql"
	"fmt"

	"batiqa-ai/internal/model"
)

// ConversationRepository handles conversations table.
type ConversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

// Create inserts a conversation entry (parameterized).
func (r *ConversationRepository) Create(c *model.Conversation) error {
	if !model.IsValidRole(c.Role) {
		return fmt.Errorf("invalid role: %s", c.Role)
	}
	query := `INSERT INTO conversations (session_id, role, message, intent) VALUES (?, ?, ?, ?)`
	res, err := r.db.Exec(query, c.SessionID, c.Role, c.Message, c.Intent)
	if err != nil {
		return fmt.Errorf("conversation Create: %w", err)
	}
	id, _ := res.LastInsertId()
	c.ID = id
	return nil
}

// ListBySession returns conversations for a session ordered by created_at asc.
func (r *ConversationRepository) ListBySession(sessionID string, limit int) ([]*model.Conversation, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	query := `SELECT id, session_id, role, message, intent, created_at FROM conversations WHERE session_id = ? ORDER BY created_at ASC LIMIT ?`
	rows, err := r.db.Query(query, sessionID, limit)
	if err != nil {
		return nil, fmt.Errorf("conversation ListBySession: %w", err)
	}
	defer rows.Close()

	var out []*model.Conversation
	for rows.Next() {
		var c model.Conversation
		var intent sql.NullString
		if err := rows.Scan(&c.ID, &c.SessionID, &c.Role, &c.Message, &intent, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("conversation scan: %w", err)
		}
		if intent.Valid {
			c.Intent = &intent.String
		}
		out = append(out, &c)
	}
	return out, rows.Err()
}
