package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"batiqa-ai/internal/config"
	"batiqa-ai/internal/router"
)

func main() {
	cfg := config.Load()

	// Try to open DB for Phase 4 ticket system. If fails, start with health only (graceful degradation per ERROR FLOW.md)
	var r http.Handler
	db, err := config.OpenDB(cfg)
	if err != nil {
		log.Printf("WARNING: DB not available (%v) - starting with health endpoint only", err)
		r = router.New()
	} else {
		log.Printf("Database connected")
		defer db.Close()
		r = router.NewWithDB(db)
	}

	srv := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background
	go func() {
		log.Printf("BATIQA AI Guest Assistant starting on %s (env=%s)", cfg.Addr(), cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
