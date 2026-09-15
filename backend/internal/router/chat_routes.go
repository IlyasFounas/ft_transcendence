package router

import (
	// "backend/internal/websockets/client"
	"backend/internal/websockets/chat"
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

func registerChatRoutes(r *gin.Engine, wsChat *chat_ws.Chat) {
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.GET("/ws/chat", chat_ws.JoinChat(wsChat))
	}
}