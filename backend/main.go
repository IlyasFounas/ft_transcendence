package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
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
    r.GET("/", func(c *gin.Context) {
        c.String(200, "Hello from Gin with hot reload")
    })
    r.Run(":8080")
}