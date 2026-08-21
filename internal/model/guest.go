package model

import "time"

// Guest represents guests table per docs/DATABASE.md
type Guest struct {
	ID         int64     `json:"id" db:"id"`
	SessionID  string    `json:"session_id" db:"session_id"`
	RoomNumber *string   `json:"room_number,omitempty" db:"room_number"`
	Language   string    `json:"language" db:"language"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
