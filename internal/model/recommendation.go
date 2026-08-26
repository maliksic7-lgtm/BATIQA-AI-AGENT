package model

import "time"

// Recommendation represents the recommendations collection
type Recommendation struct {
	ID          int64     `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Category    string    `json:"category" bson:"category"`
	Description *string   `json:"description,omitempty" bson:"description,omitempty"`
	PriceMin    *int      `json:"price_min,omitempty" bson:"price_min,omitempty"`
	PriceMax    *int      `json:"price_max,omitempty" bson:"price_max,omitempty"`
	DistanceKm  *float64  `json:"distance_km,omitempty" bson:"distance_km,omitempty"`
	Address     *string   `json:"address,omitempty" bson:"address,omitempty"`
	Active      bool      `json:"active" bson:"active"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}
