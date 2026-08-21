package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// OpenDB creates a MySQL connection pool using database/sql.
// It validates DSN, sets pool limits, and pings with timeout.
func OpenDB(cfg *Config) (*sql.DB, error) {
	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is empty")
	}
	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	db.SetMaxOpenConns(cfg.DBMaxOpen)
	db.SetMaxIdleConns(cfg.DBMaxIdle)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping with timeout (5s)
	done := make(chan error, 1)
	go func() { done <- db.Ping() }()

	select {
	case err := <-done:
		if err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("db ping failed: %w", err)
		}
	case <-time.After(5 * time.Second):
		_ = db.Close()
		return nil, fmt.Errorf("db ping timeout (5s)")
	}

	return db, nil
}
