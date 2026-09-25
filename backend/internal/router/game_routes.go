package router

import (
	"backend/internal/handlers"
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

func registerGameRoutes(r *gin.Engine, h *handlers.GameHandler) {
	public := r.Group("/")
	public.POST("/game/create", h.CreateGame())
	
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.GET("/game/:id", h.Game())
		private.GET("/game/:id/title", h.Title())
		private.GET("/game/:id/mode", h.Mode())
		private.GET("/game/:id/difficulty", h.Difficulty())
		private.GET("/game/:id/timer", h.Timer())
		private.GET("/game/:id/roundnb", h.RoundTotal())
		private.GET("/game/:id/finished", h.Finished())
		private.GET("/game/:id/playernb", h.PlayerNb())
		private.GET("/game/:id/createdat", h.CreatedAt())

		private.PATCH("/game/:id/title", h.SetTitle())
		private.PATCH("/game/:id/mode", h.SetMode())
		private.PATCH("/game/:id/difficulty", h.SetDifficulty())
		private.PATCH("/game/:id/finished", h.SetFinished())
		private.PATCH("/game/:id/playernb", h.SetPlayerNb())
	}
}