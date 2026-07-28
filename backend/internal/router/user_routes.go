package router

import (
    "github.com/uptrace/bun"
    "github.com/gin-gonic/gin"
    
    // OUR PACKAGE
    "backend/internal/handlers"
)

func RegisterUserPublicRoutes(group *gin.RouterGroup, db *bun.DB) {
    group.POST("/user/create", handler.CreateUser(db))
    group.GET("/user/activation", handler.ActivationToken(db))
}

func RegisterUserPrivateRoutes(group *gin.RouterGroup, db *bun.DB) {
	group.GET("/user/get", handler.GetUser(db))
	group.GET("/user/get/name", handler.GetName(db))
	group.GET("/user/get/login", handler.GetLogin(db))
	group.GET("/user/get/email", handler.GetEmail(db))
	group.GET("/user/get/oauth", handler.GetOAuth(db))


	group.PATCH("/user/name", handler.SetName(db))
	group.PATCH("/user/login", handler.SetLogin(db))
	group.PATCH("/user/avatar", handler.SetAvatar(db))
}
