package router

import (
	"backend/internal/handlers"
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

func registerLobbyRoutes(r *gin.Engine, h *handlers.GameHandler) {
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.POST("/lobby/create", h.CreateLobby())
	}
}