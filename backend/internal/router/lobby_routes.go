package router

import (
	"sync"
	"backend/internal/handlers"
	"backend/internal/websockets/lobby"
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

func registerLobbyRoutes(r *gin.Engine, h *handlers.GameHandler,
	hl *handlers.LobbyHandler, hub *lobby_ws.Hub, wg *sync.WaitGroup) {
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.POST("/lobby/create", h.CreateLobby())
		private.POST("/lobby/join", hl.JoinLobby())
		private.GET("/ws/lobby/join", lobby_ws.JoinLobby(hub, wg))
		private.DELETE("/lobby/leave", hl.LeaveLobby())
		private.DELETE("/lobby/remove/:username", hl.RemovePlayer())
		private.GET("/game/historic/:userid", hl.Historic())
		private.GET("/game/nbgames/:userid", hl.NbGames())
		private.GET("/lobby/allPlayers", hl.AllPlayers())
	}
}