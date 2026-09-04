package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// loginAttemptsCol stores one document per failed staff login attempt. Each
// document carries a TTL index on `ts` so Mongo auto-purges attempts older than
// the lockout window — no manual GC needed, and counts survive server restarts.
const loginAttemptsCol = "login_attempts"

// LoginRateLimitRepository backs a persistent, multi-instance rate limiter for
// staff login (vs the in-memory default used by isolated unit tests).
type LoginRateLimitRepository struct {
	db *mongo.Database
}

func NewLoginRateLimitRepository(db *mongo.Database) *LoginRateLimitRepository {
	return &LoginRateLimitRepository{db: db}
}

// Record inserts a single failed attempt keyed by (key, ts).
func (r *LoginRateLimitRepository) Record(key string, ts time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := r.db.Collection(loginAttemptsCol).InsertOne(ctx, bson.M{
		"key": key,
		"ts":  ts,
	})
	return err
}

// CountRecent returns the number of failed attempts for key since `since`.
func (r *LoginRateLimitRepository) CountRecent(key string, since time.Time) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return r.db.Collection(loginAttemptsCol).CountDocuments(ctx, bson.M{
		"key": key,
		"ts":  bson.M{"$gte": since},
	})
}

// Clear removes all recorded attempts for key (called on successful login).
func (r *LoginRateLimitRepository) Clear(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, _ = r.db.Collection(loginAttemptsCol).DeleteMany(ctx, bson.M{"key": key})
}
