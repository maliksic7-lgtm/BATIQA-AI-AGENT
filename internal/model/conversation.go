package model

import "time"

// ConversationRole enum per spec
const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

// Conversation represents conversations table
type Conversation struct {
	ID        int64     `json:"id" db:"id"`
	SessionID string    `json:"session_id" db:"session_id"`
	Role      string    `json:"role" db:"role"`
	Message   string    `json:"message" db:"message"`
	Intent    *string   `json:"intent,omitempty" db:"intent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// IsValidRole validates role value
func IsValidRole(r string) bool {
	return r == RoleUser || r == RoleAssistant || r == RoleSystem
}
