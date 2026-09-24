package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"backend/internal/handlers"
	"backend/internal/websockets/lobby"
	"backend/internal/websockets/chat"
)

func NewRouter(h *handlers.Handlers, hub *lobby_ws.Hub, 
	chat *chat_ws.Chat) *gin.Engine {
	r := gin.Default()    //when you first create a var :=
	r.Use(cors.Default()) // TODO: global middleware now but accept only localhost:3000 request

	// Splitting routes
	registerUserRoutes(r, h.User, h.Token)
	registerLocationRoutes(r, h.Location)
	registerGameRoutes(r, h.Game)
	registerLobbyRoutes(r, h.Game, h.Lobby, hub)
	registerRoundRoutes(r, h.Round)
	registerGuessRoutes(r, h.Guess)
	registerFriendRoutes(r, h.Friend)
	registerChatRoutes(r, chat)
	return r
}
