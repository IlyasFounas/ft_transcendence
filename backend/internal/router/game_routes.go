package router

import (
	"backend/internal/handlers"
	"github.com/gin-gonic/gin"
	// "backend/internal/middleware"
)

func registerGameRoutes(r *gin.Engine, h *handlers.GameHandler) {
	private := r.Group("/")
	// private.use(middleware.AuthMiddleWare())
	{
		private.POST("/game/create", h.CreateGame())
		private.GET("/game/:id", h.GetGame())
		private.GET("/game/:id/title", h.GetTitle())
		private.GET("/game/:id/mode", h.GetMode())
		private.GET("/game/:id/timer", h.GetTimer())
		private.GET("/game/:id/roundnb", h.GetRoundNb())
		private.GET("/game/:id/finished", h.GetFinished())
		private.GET("/game/:id/playernb", h.GetPlayerNb())
		private.GET("/game/:id/createdat", h.GetCreatedAt())
	}
}
