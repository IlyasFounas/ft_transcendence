package router

import (
	"sync"
	"backend/internal/websockets/chat"
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

func registerChatRoutes(r *gin.Engine, chat *chat_ws.Chat, wg *sync.WaitGroup) {
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.GET("/ws/chat", chat_ws.JoinChat(chat, wg))
	}
}