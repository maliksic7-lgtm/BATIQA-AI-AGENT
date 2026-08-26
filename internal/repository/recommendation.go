package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"batiqa-ai/internal/model"
)

// RecommendationRepository handles the recommendations collection.
type RecommendationRepository struct {
	db *mongo.Database
}

func NewRecommendationRepository(db *mongo.Database) *RecommendationRepository {
	return &RecommendationRepository{db: db}
}

const recommendationsCol = "recommendations"

// ListActive returns active recommendations with optional filters.
func (r *RecommendationRepository) ListActive(category *string, maxPrice *int) ([]*model.Recommendation, error) {
	filter := bson.M{"active": true}
	if category != nil && *category != "" {
		filter["category"] = *category
	}
	if maxPrice != nil {
		// Items without a price cap always match; otherwise price_max <= budget
		filter["$or"] = []bson.M{
			{"price_max": bson.M{"$exists": false}},
			{"price_max": nil},
			{"price_max": bson.M{"$lte": *maxPrice}},
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.db.Collection(recommendationsCol).Find(
		ctx,
		filter,
		options.Find().SetSort(bson.M{"distance_km": 1}),
	)
	if err != nil {
		return nil, fmt.Errorf("recommendation ListActive: %w", err)
	}
	defer cursor.Close(ctx)

	var out []*model.Recommendation
	if err := cursor.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("recommendation decode: %w", err)
	}
	return out, nil
}
