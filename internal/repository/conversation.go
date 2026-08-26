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

// ConversationRepository handles the conversations collection.
type ConversationRepository struct {
	db *mongo.Database
}

func NewConversationRepository(db *mongo.Database) *ConversationRepository {
	return &ConversationRepository{db: db}
}

const conversationsCol = "conversations"

// Create inserts a conversation entry.
func (r *ConversationRepository) Create(c *model.Conversation) error {
	if !model.IsValidRole(c.Role) {
		return fmt.Errorf("invalid role: %s", c.Role)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := nextID(ctx, r.db, conversationsCol)
	if err != nil {
		return fmt.Errorf("conversation Create id: %w", err)
	}
	c.ID = id
	c.CreatedAt = time.Now()
	_, err = r.db.Collection(conversationsCol).InsertOne(ctx, c)
	if err != nil {
		return fmt.Errorf("conversation Create: %w", err)
	}
	return nil
}

// ListBySession returns conversations for a session ordered by created_at asc.
func (r *ConversationRepository) ListBySession(sessionID string, limit int) ([]*model.Conversation, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.db.Collection(conversationsCol).Find(
		ctx,
		bson.M{"session_id": sessionID},
		options.Find().SetSort(bson.M{"created_at": 1}).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, fmt.Errorf("conversation ListBySession: %w", err)
	}
	defer cursor.Close(ctx)

	var out []*model.Conversation
	if err := cursor.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("conversation decode: %w", err)
	}
	return out, nil
}
