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
    db := InitDatabase();

    r.Use(cors.Default()) // TODO accept only localhost:3000 request
    var myTruc *Truc = &Truc{Title: "Hello"}

    ft_strcpy(myTruc, "Hello world from the backend!!")
    r.GET("/", func(c *gin.Context) {
        c.String(200, myTruc.Title)
    })
    r.Run(":8080")
}
