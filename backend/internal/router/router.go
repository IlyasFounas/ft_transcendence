package router

import (
	"backend/internal/handlers"
	"backend/internal/websockets/chat"
	"backend/internal/websockets/lobby"
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(ctx context.Context, h *handlers.Handlers, hub *lobby_ws.Hub,
	chat *chat_ws.Chat, wg *sync.WaitGroup, acceptFriend chan *int64) *gin.Engine {
	r := gin.Default() //when you first create a var :=
	if os.Getenv("ENV") == "prod" {
		config := cors.Config{
			AllowOrigins:     []string{os.Getenv("PUBLIC_URL")},
			AllowCredentials: true,
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodDelete,
				http.MethodPatch,
				http.MethodPut}}
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
	registerFriendRoutes(r, h.Friend, acceptFriend)
	registerChatRoutes(ctx, r, chat, wg)
	return r
}
