package main

import (
    "os"
    "context"
    "backend/internal/database"
    "backend/internal/router"
    "github.com/gin-gonic/gin"
)

/**
 * Currently, this is a basic main function to test the database connection and
 * run a basic server on port 8080.
 * The `exit` function should only be called from `main()` because if called elsewhere,
 * the garbage collector may not properly free all resources.
 */
func main() {
    env := os.Getenv("ENV")
    if (len(env) != 0 && env == "prod") {
        gin.SetMode(gin.ReleaseMode)
    }
    db := database.InitDatabase()
    
    if err := db.PingContext(context.Background()); err != nil {
        os.Exit(1)
    }
    r, err := router.InitRouter(db)
    if err != nil {
        os.Exit(1)
    }
    r.Run(":8080")
}
