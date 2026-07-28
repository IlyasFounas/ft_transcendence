package router

import (
    "github.com/uptrace/bun"
    "github.com/gin-gonic/gin"
    
    // OUR PACKAGE
    "backend/internal/handlers"
)

func RegisterAuthPublicRoutes(group *gin.RouterGroup, db *bun.DB) {
    group.POST("/auth/login", handler.LoginUser(db))
    group.GET("/auth/refresh", handler.RefreshToken(db))
}
