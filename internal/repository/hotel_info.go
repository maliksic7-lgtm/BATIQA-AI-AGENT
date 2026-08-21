package repository

import (
	"database/sql"
	"fmt"

	"batiqa-ai/internal/model"
)

// HotelInfoRepository handles hotel_information table.
type HotelInfoRepository struct {
	db *sql.DB
}

func NewHotelInfoRepository(db *sql.DB) *HotelInfoRepository {
	return &HotelInfoRepository{db: db}
}

// ListActive returns all active hotel information, optionally filtered by category.
func (r *HotelInfoRepository) ListActive(category *string) ([]*model.HotelInformation, error) {
	base := `SELECT id, category, title, content, active, created_at, updated_at FROM hotel_information WHERE active = TRUE`
	args := []interface{}{}
	if category != nil && *category != "" {
		base += ` AND category = ?`
		args = append(args, *category)
	}
	base += ` ORDER BY category ASC, title ASC`

	rows, err := r.db.Query(base, args...)
	if err != nil {
		return nil, fmt.Errorf("hotel_info ListActive: %w", err)
	}
	defer rows.Close()

	var out []*model.HotelInformation
	for rows.Next() {
		var h model.HotelInformation
		if err := rows.Scan(&h.ID, &h.Category, &h.Title, &h.Content, &h.Active, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, fmt.Errorf("hotel_info scan: %w", err)
		}
		out = append(out, &h)
	}
	return out, rows.Err()
}

// FindByID returns hotel info by id.
func (r *HotelInfoRepository) FindByID(id int64) (*model.HotelInformation, error) {
	query := `SELECT id, category, title, content, active, created_at, updated_at FROM hotel_information WHERE id = ?`
	row := r.db.QueryRow(query, id)
	var h model.HotelInformation
	err := row.Scan(&h.ID, &h.Category, &h.Title, &h.Content, &h.Active, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("hotel_info FindByID: %w", err)
	}
	return &h, nil
}

// Create inserts new hotel information (admin only, not AI).
func (r *HotelInfoRepository) Create(h *model.HotelInformation) error {
	query := `INSERT INTO hotel_information (category, title, content, active) VALUES (?, ?, ?, ?)`
	res, err := r.db.Exec(query, h.Category, h.Title, h.Content, h.Active)
	if err != nil {
		return fmt.Errorf("hotel_info Create: %w", err)
	}
	id, _ := res.LastInsertId()
	h.ID = id
	return nil
}
