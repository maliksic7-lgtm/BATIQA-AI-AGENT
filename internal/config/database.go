package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo connects to the MongoDB server and returns the database handle.
// Returns an error if the server cannot be reached within 5 seconds,
// allowing graceful degradation per ERROR FLOW.md.
func ConnectMongo(cfg *Config) (*mongo.Database, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, nil, err
	}

	db := client.Database(cfg.MongoDB)
	cleanup := func() { _ = client.Disconnect(context.Background()) }
	return db, cleanup, nil
}
