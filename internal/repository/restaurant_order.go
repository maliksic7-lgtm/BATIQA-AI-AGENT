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

// restaurantOrdersCol stores restaurant room-service orders placed via chat.
const restaurantOrdersCol = "restaurant_orders"

// RestaurantOrderRepository handles CRUD + analytics for restaurant orders.
type RestaurantOrderRepository struct {
	db *mongo.Database
}

func NewRestaurantOrderRepository(db *mongo.Database) *RestaurantOrderRepository {
	return &RestaurantOrderRepository{db: db}
}

// Create persists a new order. It atomically allocates the next sequence id and
// derives the human-friendly ORD-XXXXX order number.
func (r *RestaurantOrderRepository) Create(o *model.RestaurantOrder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := nextID(ctx, r.db, "restaurant_order")
	if err != nil {
		return err
	}
	o.ID = id
	o.OrderNumber = fmt.Sprintf("ORD-%05d", id)
	total := 0
	for _, it := range o.Items {
		total += it.Price * it.Quantity
	}
	o.TotalPrice = total
	o.Status = model.OrderNew
	o.CreatedAt = time.Now()
	o.UpdatedAt = o.CreatedAt
	_, err = r.db.Collection(restaurantOrdersCol).InsertOne(ctx, o)
	return err
}

// ListBySession returns orders for a guest session, newest first (capped 50).
func (r *RestaurantOrderRepository) ListBySession(sessionID string, limit int) ([]model.RestaurantOrder, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.db.Collection(restaurantOrdersCol).Find(
		ctx,
		bson.M{"session_id": sessionID},
		options.Find().SetSort(bson.M{"_id": -1}).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, fmt.Errorf("restaurant ListBySession: %w", err)
	}
	defer cur.Close(ctx)
	out := []model.RestaurantOrder{}
	for cur.Next(ctx) {
		var o model.RestaurantOrder
		if err := cur.Decode(&o); err == nil {
			out = append(out, o)
		}
	}
	return out, nil
}

// ListAll returns recent orders across all rooms (for staff dashboard), newest
// first, with optional status filter, capped at `limit`.
func (r *RestaurantOrderRepository) ListAll(status string, limit int) ([]model.RestaurantOrder, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := bson.M{}
	if status != "" {
		query["status"] = status
	}
	cur, err := r.db.Collection(restaurantOrdersCol).Find(
		ctx,
		query,
		options.Find().SetSort(bson.M{"_id": -1}).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, fmt.Errorf("restaurant ListAll: %w", err)
	}
	defer cur.Close(ctx)
	out := []model.RestaurantOrder{}
	for cur.Next(ctx) {
		var o model.RestaurantOrder
		if err := cur.Decode(&o); err == nil {
			out = append(out, o)
		}
	}
	return out, nil
}

// UpdateStatus transitions an order's lifecycle status (staff action).
func (r *RestaurantOrderRepository) UpdateStatus(id int64, status string) (model.RestaurantOrder, bool, error) {
	if !model.IsValidOrderStatus(status) {
		return model.RestaurantOrder{}, false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var o model.RestaurantOrder
	err := r.db.Collection(restaurantOrdersCol).FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&o)
	if err != nil {
		return model.RestaurantOrder{}, false, err
	}
	return o, true, nil
}

// TopOrderedItems aggregates line-item quantities across placed restaurant
// orders (not just chat mentions), returning the most-ordered F&B by total
// quantity. Returns CategoryCount{Category: item-name, Count: total quantity}.
func (r *RestaurantOrderRepository) TopOrderedItems(limit int) ([]CategoryCount, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.db.Collection(restaurantOrdersCol).Aggregate(ctx, []bson.M{
		{"$match": bson.M{"status": bson.M{"$ne": model.OrderCancelled}}},
		{"$unwind": "$items"},
		{"$group": bson.M{"_id": "$items.name", "qty": bson.M{"$sum": "$items.quantity"}}},
		{"$sort": bson.M{"qty": -1}},
		{"$limit": limit},
	})
	if err != nil {
		return nil, fmt.Errorf("restaurant TopOrderedItems: %w", err)
	}
	defer cur.Close(ctx)
	out := []CategoryCount{}
	for cur.Next(ctx) {
		var row struct {
			ID  string `bson:"_id"`
			Qty int    `bson:"qty"`
		}
		if err := cur.Decode(&row); err == nil {
			out = append(out, CategoryCount{Category: row.ID, Count: int64(row.Qty)})
		}
	}
	return out, nil
}
