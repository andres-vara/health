package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andres-vara/health"
)

func main() {
	// Create a new router
	router := http.NewServeMux()

	// Add health endpoints - one with plain text and one with JSON
	router.Handle("/health", health.Handler())
	router.Handle("/health/json", health.Handler().WithJSON(true))

	// Add a route to toggle health status for demonstration
	router.HandleFunc("/toggle-health", func(w http.ResponseWriter, r *http.Request) {
		if health.GetStatus() == health.Up {
			health.SetUnhealthy("Service manually marked as unhealthy")
			fmt.Fprintf(w, "Health status set to DOWN with reason")
		} else {
			health.SetHealthy()
			fmt.Fprintf(w, "Health status set to UP")
		}
	})

	// Start server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Print usage instructions
	log.Println("Server started. Available endpoints:")
	log.Println("- GET /health - Health check endpoint (plain text)")
	log.Println("- GET /health/json - Health check endpoint (JSON)")
	log.Println("- GET /toggle-health - Toggle health status for demonstration")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	
	// Simulate unhealthy state during shutdown
	health.SetUnhealthy("Server is shutting down")
	
	// Give existing connections time to complete
	time.Sleep(time.Second)
	
	log.Println("Server stopped")
} 