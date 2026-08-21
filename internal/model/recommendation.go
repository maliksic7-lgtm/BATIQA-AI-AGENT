package model

import "time"

// Recommendation represents recommendations table
type Recommendation struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Category    string    `json:"category" db:"category"`
	Description *string   `json:"description,omitempty" db:"description"`
	PriceMin    *int      `json:"price_min,omitempty" db:"price_min"`
	PriceMax    *int      `json:"price_max,omitempty" db:"price_max"`
	DistanceKm  *float64  `json:"distance_km,omitempty" db:"distance_km"`
	Address     *string   `json:"address,omitempty" db:"address"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
