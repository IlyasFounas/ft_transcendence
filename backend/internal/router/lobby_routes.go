package router

import (
	"backend/internal/handlers"
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

func registerLobbyRoutes(r *gin.Engine, h *handlers.GameHandler,
	hl *handlers.LobbyHandler) {
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.POST("/ws/lobby/create", h.CreateLobby())
		private.POST("/ws/lobby/join", hl.JoinLobby())
		private.DELETE("/lobby/leave", hl.LeaveLobby())
		private.DELETE("/lobby/remove/:username", hl.RemovePlayer())
		private.GET("/lobby/allPlayers", hl.AllPlayers())
	}
}