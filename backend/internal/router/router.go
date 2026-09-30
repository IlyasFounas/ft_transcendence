package router

import (
	"context"
	"sync"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"backend/internal/handlers"
	"backend/internal/websockets/lobby"
	"backend/internal/websockets/chat"
)

func NewRouter(ctx context.Context, h *handlers.Handlers, hub *lobby_ws.Hub, 
	chat *chat_ws.Chat, wg *sync.WaitGroup) *gin.Engine {
	r := gin.Default()    //when you first create a var :=
	r.Use(cors.Default()) // TODO: global middleware now but accept only localhost:3000 request

	// Splitting routes
	registerUserRoutes(r, h.User, h.Token)
	registerStatisticsRoutes(r, h.Statistics)
	registerLocationRoutes(r, h.Location)
	registerGameRoutes(r, h.Game)
	registerLobbyRoutes(ctx, r, h.Game, h.Lobby, hub, wg)
	registerGuessRoutes(r, h.Guess)
	registerFriendRoutes(r, h.Friend)
	registerChatRoutes(ctx, r, chat, wg)
	return r
}