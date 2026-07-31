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
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/router"
)

/**
 * Currently, this is a basic main function to test the database connection and
 * run a basic server on port 8080.
 * The `exit` function should only be called from `main()` because if called elsewhere,
 * the garbage collector may not properly free all resources.
 */
func main() {
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
	locationRepo, err := repository.NewLocationRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize location repo: %v", err)
	}

	// 3. Instanciation of handlers
	handlers := handler.NewHandler(locationRepo)

	// 4. Router setup
	r := router.NewRouter(handlers)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 5. Run the server in a goroutine to allow graceful shutdown
	go func() {
		log.Println("Starting server on :8080...")
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
	log.Println("Closing repository and flushing logs...")
	locationRepo.Close()

	log.Println("Application closed properly.")
}
