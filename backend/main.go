package main

import (
	// EXTERNAL PACKAGE
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// OUR PACKAGE
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/repositories"
	"backend/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/sams96/rgeo"
)

/**
 * Initializes the database, repositories, handlers and router, then starts
 * the HTTP server on :8080 in a goroutine. Waits for SIGINT/SIGTERM to
 * trigger a graceful shutdown (server drains active requests, then the
 * repository is closed). Fatal errors use log.Fatalf, which calls
 * os.Exit(1) internally — so it must only be called from main(), otherwise
 * deferred cleanup and graceful shutdown would be skipped.
 */
func main() {
    env := os.Getenv("ENV")
    if (len(env) != 0 && env == "prod") {
        gin.SetMode(gin.ReleaseMode)
    }
	// 1. Connection to the database
	db := database.InitDatabase()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error during BDD closure: %v", err)
		}
	}()
	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("Unable to reach the database: %v", err)
	}

	// 2. Instanciation of repositories
	repo, err := repositories.NewRepositories(db)
	if err != nil {
		log.Fatalf("Failed to initialize location repo: %v", err)
	}

	// 3. Instanciation of handlers
	reverser, err := rgeo.New(rgeo.Countries110)
	if err != nil {
		log.Fatalf("Failed to initialize country reverser: %v", err)
	}
	h := handlers.NewHandlers(repo, reverser)

	// 4. Router setup
	r := router.NewRouter(h)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 5. Run the server in a goroutine to allow graceful shutdown
	go func() {
		log.Println("Starting server and listening on :8080...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Unexpected error from HTTP server: %v", err)
		}
	}()

	// 6. Graceful shutdown on SIGINT or SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received. Beginning graceful shutdown...")

	// 7. Shutdown the server with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error during forced shutdown of HTTP server: %v", err)
	}

	// 8. Final cleanup of the asynchronous logs buffer
	log.Println("Closing location repository and flushing logs...")
	repo.Location.Close()
	log.Println("Application closed properly.")
}
