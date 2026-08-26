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

// GuestRepository handles the guests collection.
type GuestRepository struct {
	db *mongo.Database
}

func NewGuestRepository(db *mongo.Database) *GuestRepository {
	return &GuestRepository{db: db}
}

const guestsCol = "guests"

func guestCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// FindBySessionID returns guest by session_id.
func (r *GuestRepository) FindBySessionID(sessionID string) (*model.Guest, error) {
	ctx, cancel := guestCtx()
	defer cancel()

	var g model.Guest
	err := r.db.Collection(guestsCol).FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&g)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("guest FindBySessionID: %w", err)
	}
	return &g, nil
}

// Upsert creates or updates a guest session. Room number is never invented and
// an omitted room never overwrites an existing one; caller provides validated values.
func (r *GuestRepository) Upsert(sessionID string, roomNumber *string, language string) (*model.Guest, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}
	if language == "" {
		language = "id"
	}
	ctx, cancel := guestCtx()
	defer cancel()

	now := time.Now()
	existing, err := r.FindBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		g := &model.Guest{SessionID: sessionID, Language: language, CreatedAt: now, UpdatedAt: now}
		if roomNumber != nil && *roomNumber != "" {
			g.RoomNumber = roomNumber
		}
		id, err := nextID(ctx, r.db, guestsCol)
		if err != nil {
			return nil, fmt.Errorf("guest Upsert id: %w", err)
		}
		g.ID = id
		if _, err := r.db.Collection(guestsCol).InsertOne(ctx, g); err != nil {
			return nil, fmt.Errorf("guest Upsert insert: %w", err)
		}
		return g, nil
	}

	set := bson.M{"language": language, "updated_at": now}
	if roomNumber != nil && *roomNumber != "" {
		set["room_number"] = *roomNumber // never overwrite with empty
	}
	var g model.Guest
	err = r.db.Collection(guestsCol).FindOneAndUpdate(
		ctx,
		bson.M{"_id": existing.ID},
		bson.M{"$set": set},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&g)
	if err != nil {
		return nil, fmt.Errorf("guest Upsert update: %w", err)
	}
	return &g, nil
}

// Create inserts a new guest. Returns error if session_id already exists.
func (r *GuestRepository) Create(g *model.Guest) error {
	ctx, cancel := guestCtx()
	defer cancel()

	id, err := nextID(ctx, r.db, guestsCol)
	if err != nil {
		return fmt.Errorf("guest Create id: %w", err)
	}
	g.ID = id
	now := time.Now()
	g.CreatedAt = now
	g.UpdatedAt = now
	_, err = r.db.Collection(guestsCol).InsertOne(ctx, g)
	if err != nil {
		return fmt.Errorf("guest Create: %w", err)
	}
	return nil
}
