package router

import (
	"backend/internal/middleware"
	"backend/internal/websockets/chat"
	"context"
	"github.com/gin-gonic/gin"
	"sync"
)

func registerChatRoutes(ctx context.Context, r *gin.Engine, chat *chat_ws.Chat,
	wg *sync.WaitGroup) {
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.GET("/ws/chat", chat_ws.JoinChat(ctx, chat, wg))
	}
}
