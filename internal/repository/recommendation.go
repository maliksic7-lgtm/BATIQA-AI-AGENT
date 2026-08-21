package repository

import (
	"database/sql"
	"fmt"

	"batiqa-ai/internal/model"
)

// RecommendationRepository handles recommendations table.
type RecommendationRepository struct {
	db *sql.DB
}

func NewRecommendationRepository(db *sql.DB) *RecommendationRepository {
	return &RecommendationRepository{db: db}
}

// ListActive returns active recommendations with optional filters (parameterized).
func (r *RecommendationRepository) ListActive(category *string, maxPrice *int) ([]*model.Recommendation, error) {
	base := `SELECT id, name, category, description, price_min, price_max, distance_km, address, active, created_at FROM recommendations WHERE active = TRUE`
	args := []interface{}{}
	if category != nil && *category != "" {
		base += ` AND category = ?`
		args = append(args, *category)
	}
	if maxPrice != nil {
		base += ` AND (price_max IS NULL OR price_max <= ?)`
		args = append(args, *maxPrice)
	}
	base += ` ORDER BY distance_km ASC`

	rows, err := r.db.Query(base, args...)
	if err != nil {
		return nil, fmt.Errorf("recommendation ListActive: %w", err)
	}
	defer rows.Close()

	var out []*model.Recommendation
	for rows.Next() {
		var rec model.Recommendation
		var desc, addr sql.NullString
		var pMin, pMax sql.NullInt64
		var dist sql.NullFloat64
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.Category, &desc, &pMin, &pMax, &dist, &addr, &rec.Active, &rec.CreatedAt); err != nil {
			return nil, fmt.Errorf("recommendation scan: %w", err)
		}
		if desc.Valid {
			rec.Description = &desc.String
		}
		if pMin.Valid {
			v := int(pMin.Int64)
			rec.PriceMin = &v
		}
		if pMax.Valid {
			v := int(pMax.Int64)
			rec.PriceMax = &v
		}
		if dist.Valid {
			rec.DistanceKm = &dist.Float64
		}
		if addr.Valid {
			rec.Address = &addr.String
		}
		out = append(out, &rec)
	}
	return out, rows.Err()
}
