package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// StaffSession is a persisted staff session (survives server restarts).
// _id is sha256(token) so a DB leak never exposes raw bearer tokens.
type StaffSession struct {
	TokenHash  string    `bson:"_id"`
	StaffID    int64     `bson:"staff_id"`
	Name       string    `bson:"name"`
	Email      string    `bson:"email"`
	Department string    `bson:"department"`
	ExpiresAt  time.Time `bson:"expires_at"`
}

// StaffSessionRepository handles the staff_sessions collection.
type StaffSessionRepository struct {
	db *mongo.Database
}

func NewStaffSessionRepository(db *mongo.Database) *StaffSessionRepository {
	return &StaffSessionRepository{db: db}
}

const staffSessionsCol = "staff_sessions"

func (r *StaffSessionRepository) Create(ctx context.Context, s StaffSession) error {
	_, err := r.db.Collection(staffSessionsCol).InsertOne(ctx, s)
	return err
}

func (r *StaffSessionRepository) Find(ctx context.Context, tokenHash string) (*StaffSession, error) {
	var s StaffSession
	err := r.db.Collection(staffSessionsCol).FindOne(ctx, bson.M{"_id": tokenHash}).Decode(&s)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	if time.Now().After(s.ExpiresAt) {
		_ = r.Delete(ctx, tokenHash)
		return nil, nil
	}
	return &s, nil
}

func (r *StaffSessionRepository) Delete(ctx context.Context, tokenHash string) error {
	_, err := r.db.Collection(staffSessionsCol).DeleteOne(ctx, bson.M{"_id": tokenHash})
	return err
}
