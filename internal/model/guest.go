package model

import "time"

// Guest represents the guests collection per docs/DATABASE.md
type Guest struct {
	ID         int64     `json:"id" bson:"_id"`
	SessionID  string    `json:"session_id" bson:"session_id"`
	RoomNumber *string   `json:"room_number,omitempty" bson:"room_number,omitempty"`
	Language   string    `json:"language" bson:"language"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}
