package model

import "time"

// ConversationRole enum per spec
const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

// Conversation represents the conversations collection
type Conversation struct {
	ID        int64     `json:"id" bson:"_id"`
	SessionID string    `json:"session_id" bson:"session_id"`
	Role      string    `json:"role" bson:"role"`
	Message   string    `json:"message" bson:"message"`
	Intent    *string   `json:"intent,omitempty" bson:"intent,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// IsValidRole validates role value
func IsValidRole(r string) bool {
	return r == RoleUser || r == RoleAssistant || r == RoleSystem
}
