package main

import (
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/repositories"
	"backend/internal/router"
	chat_ws "backend/internal/websockets/chat"
	"backend/internal/websockets/handlers"
	"backend/internal/websockets/lobby"
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	_, err := url.Parse(os.Getenv("PUBLIC_URL"))
	if err != nil {
		log.Printf("error : %v", err)
		return
	}
	if os.Getenv("PUBLIC_URL") == "" {
		log.Fatalf("error : %v", err)
	}
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
	// the admin routes check in the database that the user is an admin
	middleware.SetAdminChecker(repo.User)

	// 3. Instanciation of handlers
	reverser, err := rgeo.New(rgeo.Countries110)
	if err != nil {
		log.Fatalf("Failed to initialize country reverser: %v", err)
	}
	h := handlers.NewHandlers(repo, reverser)

	// 4. Instanciation of websockets handlers
	wh := ws_handlers.NewWsHandlers(repo, reverser)

	// 5. Run the hub websocket
	signalCtx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM) 
	defer stop()
	var wg sync.WaitGroup
	wg.Add(1)
	hub := lobby_ws.NewHub(repo, wh)
	go hub.Run(signalCtx, &wg)

	wg.Add(1)
	chat := chat_ws.NewChat(repo, wh)
	go chat.Run(signalCtx, &wg)

	// 6. Router setup

	r := router.NewRouter(signalCtx, h, hub, chat, &wg)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	
	// 7. Run the server in a goroutine to allow graceful shutdown
	go func() {
		log.Println("Starting server and listening on :8080...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Unexpected error from HTTP server: %v", err)
			return
		}
	}()

	// 8. Graceful shutdown
	<-signalCtx.Done()
	
	// 9. Shutdown the server with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error during forced shutdown of HTTP server: %v", err)
	}

	// 10. Wait until every clients connexions are closed proprely
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <- done:
		log.Println("All connexion has been closed")
	case <-time.After(3*time.Second):
		log.Println("timeout occured to shutdown the server")
	}
	// 11. Final cleanup of the asynchronous logs buffer
	log.Println("Closing location repository and flushing logs...")
	repo.Location.Close()
	log.Println("Application closed properly.")
}