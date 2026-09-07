package router

import (
	"backend/internal/handlers"
	"backend/internal/websockets/lobby"
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

func registerLobbyRoutes(r *gin.Engine, h *handlers.GameHandler,
	hl *handlers.LobbyHandler) {
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.POST("/lobby/create", h.CreateLobby())
		private.POST("/lobby/join", hl.JoinLobby())
		// private.GET("/ws/lobby/join", client_ws.WsLobby())
		private.DELETE("/lobby/leave", hl.LeaveLobby())
		private.DELETE("/lobby/remove/:username", hl.RemovePlayer())
		private.GET("/lobby/allPlayers", hl.AllPlayers())
	}
	public := r.Group("/")
	public.GET("/ws/lobby/join", client_ws.WsJoinLobby())
}