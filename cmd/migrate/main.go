package main

import (
	"context"
	"log"

	"batiqa-ai/internal/config"
	"batiqa-ai/internal/repository"
)

func main() {
	cfg := config.Load()

	log.Printf("Connecting to MongoDB (%s)...", cfg.MongoDB)

	ctx := context.Background()
	db, closeDB, err := config.ConnectMongo(cfg)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	defer closeDB()

	log.Printf("Running migration + seed on database %q...", cfg.MongoDB)
	if err := repository.Migrate(ctx, db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
	for _, col := range []string{"guests", "conversations", "tickets", "hotel_information", "recommendations", "staff", "ticket_assignments"} {
		n, err := db.Collection(col).EstimatedDocumentCount(ctx)
		if err != nil {
			continue
		}
		log.Printf(" - %s: %d docs", col, n)
	}
}
