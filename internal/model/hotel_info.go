package model

import "time"

// HotelInformation represents hotel_information table
type HotelInformation struct {
	ID        int64     `json:"id" db:"id"`
	Category  string    `json:"category" db:"category"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
