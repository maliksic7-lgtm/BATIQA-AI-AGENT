package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// countersCol stores named auto-increment sequences (replaces SQL AUTO_INCREMENT).
const countersCol = "counters"

// nextID atomically increments a named counter and returns the new sequence value.
func nextID(ctx context.Context, db *mongo.Database, name string) (int64, error) {
	var result struct {
		Seq int64 `bson:"seq"`
	}
	err := db.Collection(countersCol).FindOneAndUpdate(
		ctx,
		bson.M{"_id": name},
		bson.M{"$inc": bson.M{"seq": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&result)
	if err != nil {
		return 0, fmt.Errorf("counter %s: %w", name, err)
	}
	return result.Seq, nil
}
