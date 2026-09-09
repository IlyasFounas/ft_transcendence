package router

import (
	// EXTERNAL PACKAGE
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// OUR PACKAGE
	"backend/internal/handlers"
	"backend/internal/websockets/lobby"
)

func NewRouter(h *handlers.Handlers, hub *lobby_ws.HubWs) *gin.Engine {
	r := gin.Default()    //when you first create a var :=
	r.Use(cors.Default()) // TODO: global middleware now but accept only localhost:3000 request

	// Splitting routes
	registerUserRoutes(r, h.User, h.Token)
	registerLocationRoutes(r, h.Location)
	registerGameRoutes(r, h.Game)
	registerLobbyRoutes(r, h.Game, h.Lobby, hub)
	return r
}

