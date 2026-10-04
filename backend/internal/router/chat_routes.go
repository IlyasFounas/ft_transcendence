package router

import (
	"context"
	"sync"
	"backend/internal/websockets/chat"
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

func registerChatRoutes(ctx context.Context, r *gin.Engine, chat *chat_ws.Chat,
	wg *sync.WaitGroup) {
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.GET("/ws/chat", chat_ws.JoinChat(ctx, chat, wg))
	}
}