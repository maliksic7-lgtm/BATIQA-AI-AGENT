package main

import (
	"log"
	"os"

	"batiqa-ai/internal/config"
	"batiqa-ai/internal/repository"
)

func main() {
	cfg := config.Load()

	// Allow override migrations dir via arg or env
	dir := "migrations"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	if envDir := os.Getenv("MIGRATIONS_DIR"); envDir != "" {
		dir = envDir
	}

	log.Printf("Connecting to database: %s (maxOpen=%d maxIdle=%d)", maskDSN(cfg.DBDSN), cfg.DBMaxOpen, cfg.DBMaxIdle)

	db, err := config.OpenDB(cfg)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer db.Close()

	log.Printf("Running migrations from %s...", dir)
	if err := repository.RunMigrations(db, dir); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
	// Verify tables
	rows, err := db.Query(`SHOW TABLES`)
	if err != nil {
		log.Printf("SHOW TABLES failed: %v", err)
		return
	}
	defer rows.Close()
	log.Println("Tables in database:")
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err == nil {
			log.Printf(" - %s", t)
		}
	}
}

// maskDSN hides password for logging
func maskDSN(dsn string) string {
	// dsn format: user:password@tcp(host:port)/dbname?...
	// Simple mask: replace between : and @
	at := -1
	colon := -1
	for i, c := range dsn {
		if c == ':' && colon == -1 {
			colon = i
		}
		if c == '@' {
			at = i
			break
		}
	}
	if colon != -1 && at != -1 && at > colon {
		return dsn[:colon+1] + "****" + dsn[at:]
	}
	return dsn
}
