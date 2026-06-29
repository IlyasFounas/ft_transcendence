package main

import (
    // EXTERNAL PACKAGE
    "context"
	"database/sql"
	"fmt"
    
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
    
    // OUR PACKAGE
    "internal/init_database"
)

type Truc struct{
    Title string
}

func ft_strcpy(t *Truc, to_cpy string) {
    t.Title = to_cpy
}

func main() {
    r := gin.Default() //when you first create a var :=
    db := InitDatabase();

    r.Use(cors.Default())
    // SHOULD USE THIS CODE INSTEAD OF THE DEFAULT ONE
    /* r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, // ✅ Seulement ton frontend
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true, // ✅ Si tu utilises des cookies/auth
        MaxAge:           12 * time.Hour,
    })) */
    var myTruc *Truc = &Truc{Title: "Hello"}


    ft_strcpy(myTruc, "Hello world from the backend!")
    r.GET("/", func(c *gin.Context) {
        c.String(200, myTruc.Title)
    })
    r.Run(":8080")
}