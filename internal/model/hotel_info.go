package model

import "time"

// HotelInformation represents the hotel_information collection
type HotelInformation struct {
	ID        int64     `json:"id" bson:"_id"`
	Category  string    `json:"category" bson:"category"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	Active    bool      `json:"active" bson:"active"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
