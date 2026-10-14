package router

import (
	"os"
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
	if os.Getenv("ENV") == "prod" {
		allowOrigins := []string{os.Getenv("PUBLIC_URL")}
		allowMethods := []string{"GET", "POST", "DELETE", "PATCH", "PUT"}
		config := cors.Config{
				AllowOrigins: allowOrigins,
				AllowCredentials: true,
				AllowMethods: allowMethods}
		r.Use(cors.New(config))
	} else {
		r.Use(cors.Default())
	}

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