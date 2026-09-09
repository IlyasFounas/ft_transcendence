package router

import (
	"github.com/gin-gonic/gin"
	"backend/internal/handlers"
	"backend/internal/middleware"
)

func registerFriendRoutes(r *gin.Engine, h *handlers.FriendHandler) {
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.POST("/friend/add", h.AddFriend())
		private.GET("/friend/requests", h.FriendRequests())
		private.PATCH("/friend/requests/:id/accept", h.AcceptRequest())
	}
}