package router

import (
	"github.com/gin-gonic/gin"
	"backend/internal/handlers"
	"backend/internal/middleware"
)

func registerFriendRoutes(r *gin.Engine, h *handlers.FriendHandler,
	acceptFriend chan *int64) {
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.POST("/friend/add", h.AddFriend())

		private.GET("/friend", h.Friends())
		private.GET("/friend/requests", h.FriendRequests())
		private.GET("/friend/requests/pending", h.FriendPending())
		private.GET("/friend/status/:username", h.FriendStatus())

		private.PATCH("/friend/requests/:id/accept", h.AcceptRequest(acceptFriend))

		private.DELETE("/friend/requests/:id/reject", h.RejectRequest())
		private.DELETE("/friend/requests/:id/cancel", h.CancelRequest())
		private.DELETE("/friend/:id/remove", h.RemoveFriend(acceptFriend))
	}
}