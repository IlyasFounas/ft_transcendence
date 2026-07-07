package main

import (
    // EXTERNAL PACKAGE
    "os"
    "context"
    "net/http"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    
    // OUR PACKAGE
    "backend/internal/database"
    "backend/internal/handlers"
)

/**
 * Currently, this is a basic main function to test the database connection and
 * run a basic server on port 8080.
 * The `exit` function should only be called from `main()` because if called elsewhere,
 * the garbage collector may not properly free all resources.
 */
func main() {
    r := gin.Default() //when you first create a var :=
    db := database.InitDatabase();

    r.Use(cors.Default()) // TODO accept only localhost:3000 request
    if err := db.PingContext(context.Background()); err != nil {
        os.Exit(1)
    }
    user.UserRoute(r)
    user.GameRoute(r)
    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "hello world from the backend !")
    })
    r.Run(":8080")
}