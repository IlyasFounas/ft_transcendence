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

type Truc struct{
    Title string
}

func ft_strcpy(t *Truc, to_cpy string) {
    t.Title = to_cpy
}

func main() {
    r := gin.Default() //when you first create a var :=
    db := database.InitDatabase();

    r.Use(cors.Default()) // TODO accept only localhost:3000 request
    var myTruc *Truc = &Truc{Title: "Hello"}

    if err := db.PingContext(context.Background()); err != nil {
        os.Exit(1)
    }
    ft_strcpy(myTruc, "Hello world from the backend")
    r.GET("/", func(c *gin.Context) {
        c.String(200, myTruc.Title)
    })
    r.Run(":8080")
}
