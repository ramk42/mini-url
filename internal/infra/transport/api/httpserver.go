package api

import (
	"context"
	"errors"
	"github.com/ramk42/mini-url/internal/infra/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ListenAndServe() {
	// Create Chi router
	r := createRouter()

	// Configure HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Channel to capture shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start server in goroutine
	go func() {
		log.Printf("Server started on http://localhost%s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("\nShutdown signal received - Starting graceful shutdown...")

	// Create shutdown context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Graceful shutdown error: %v", err)
	} else {
		log.Println("Server shut down cleanly")
	}
	log.Println("Closing database connections...")
	if err := database.Close(); err != nil {
		log.Printf("Graceful shutdown error: %v", err)
	}
}
