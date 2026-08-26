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

// HotelInfoRepository handles the hotel_information collection.
type HotelInfoRepository struct {
	db *mongo.Database
}

func NewHotelInfoRepository(db *mongo.Database) *HotelInfoRepository {
	return &HotelInfoRepository{db: db}
}

const hotelInfoCol = "hotel_information"

// ListActive returns all active hotel information, optionally filtered by category.
func (r *HotelInfoRepository) ListActive(category *string) ([]*model.HotelInformation, error) {
	filter := bson.M{"active": true}
	if category != nil && *category != "" {
		filter["category"] = *category
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.db.Collection(hotelInfoCol).Find(
		ctx,
		filter,
		options.Find().SetSort(bson.D{{Key: "category", Value: 1}, {Key: "title", Value: 1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("hotel_info ListActive: %w", err)
	}
	defer cursor.Close(ctx)

	var out []*model.HotelInformation
	if err := cursor.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("hotel_info decode: %w", err)
	}
	return out, nil
}

// FindByID returns hotel info by id.
func (r *HotelInfoRepository) FindByID(id int64) (*model.HotelInformation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var h model.HotelInformation
	err := r.db.Collection(hotelInfoCol).FindOne(ctx, bson.M{"_id": id}).Decode(&h)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("hotel_info FindByID: %w", err)
	}
	return &h, nil
}

// Create inserts new hotel information (admin only, not AI).
func (r *HotelInfoRepository) Create(h *model.HotelInformation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := nextID(ctx, r.db, hotelInfoCol)
	if err != nil {
		return fmt.Errorf("hotel_info Create id: %w", err)
	}
	h.ID = id
	now := time.Now()
	h.CreatedAt = now
	h.UpdatedAt = now
	_, err = r.db.Collection(hotelInfoCol).InsertOne(ctx, h)
	if err != nil {
		return fmt.Errorf("hotel_info Create: %w", err)
	}
	return nil
}
