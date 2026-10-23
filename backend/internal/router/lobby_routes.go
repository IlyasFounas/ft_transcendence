package router

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/websockets/lobby"
	"context"
	"github.com/gin-gonic/gin"
	"sync"
)

func registerLobbyRoutes(ctx context.Context, r *gin.Engine, h *handlers.GameHandler,
	hl *handlers.LobbyHandler, hub *lobby_ws.Hub, wg *sync.WaitGroup) {
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.POST("/lobby/create", h.CreateLobby())
		private.POST("/lobby/join", hl.JoinLobby())
		private.GET("/ws/lobby/join", lobby_ws.JoinLobby(ctx, hub, wg))
		private.GET("/lobby/current", hl.CurrentLobby())
		private.DELETE("/lobby/leave", hl.LeaveLobby())
		private.DELETE("/lobby/remove/:username", hl.RemovePlayer())
		private.GET("/game/historic/:userid", hl.Historic())
		private.GET("/game/nbgames/:userid", hl.NbGames())
		private.GET("/lobby/allPlayers", hl.AllPlayers())
	}
}
