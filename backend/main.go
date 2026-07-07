package main

import (
    // EXTERNAL PACKAGE
    "os"
    "context"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    
    // OUR PACKAGE
    "backend/internal/database"
)

// Similar to C, we can define structs
type Truc struct{
    Title string
}

// We can use pointers to modify data more easily
func ft_strcpy(t *Truc, to_cpy string) {
    t.Title = to_cpy
}

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
    var myTruc *Truc = &Truc{Title: "Hello"}

    if err := db.PingContext(context.Background()); err != nil {
        os.Exit(1)
    }
    ft_strcpy(myTruc, "Hello world from the backend!")
    r.GET("/", func(c *gin.Context) {
        c.String(200, myTruc.Title)
    })
    r.Run(":8080")
}
