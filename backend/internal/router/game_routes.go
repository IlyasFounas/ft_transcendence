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
	}
}
