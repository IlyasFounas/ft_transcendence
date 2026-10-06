package router

import (
	// EXTERNAL PACKAGE
	"github.com/gin-gonic/gin"

	// OUR PACKAGE
	"backend/internal/handlers"
	"backend/internal/middleware"
)

/**
 * The router and all its endpoints are now split up from main.go
 * for better clarity and maintainability.
 *
 * When working on the frontend, refer to this file to check which functions
 * are available and what endpoints they expose (e.g., /user/login to authenticate users).
 *
 * Note that endpoints under the private group are only accessible
 * when the user is authenticated.
 *
 * The difference between PUT and PATCH in the database:
 * PUT = modify a lot of field
 * PATCH = modify a unique field
 */
func registerUserRoutes(r *gin.Engine, h *handlers.UserHandler, ht *handlers.TokenHandler) {
	r.Static("/uploads", "./uploads")
	public := r.Group("/")
	{
		public.POST("/user/create", h.CreateUser())
		public.POST("/auth/login", h.Login())
		public.GET("/auth/google", h.GoogleLogin())
		public.GET("/auth/google/callback", h.GoogleCallback())
		// BUGFIX: ces deux routes appelaient h.ActivationToken() / h.RefreshToken(),
		// des méthodes qui n'existent pas sur UserHandler — elles vivent sur TokenHandler.
		public.GET("/user/activation", ht.ActivationToken())
		public.GET("/auth/refresh", ht.RefreshToken())
		public.GET("/leaderboard", h.Leaderboard()) // 100 first players ranked by xp
	}
	// it's private because only a registred user can access it
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.GET("/user", h.User())
		private.GET("/user/name", h.Name())
		private.GET("/user/login", h.Login())
		private.GET("/user/email", h.Email())
		// private.GET("/user/:username", h.UserByName())
		private.PATCH("/user/name", h.SetName())
		private.PATCH("/user/login", h.SetLogin())
		private.PATCH("/user/avatar", h.SetAvatar())
		private.POST("/auth/logout", h.Logout())
		private.POST("/user/avatar/upload", h.UploadAvatar())
		private.GET("/user/xp/:username", h.XPByName())
		private.GET("/leaderboard/user", h.UserLeaderboard()) // position, name and xp of the connected user
	}
	// TODO admin: to be enabled in the admin group (a user must not be able to set his own xp)
	// admin.PATCH("/admin/user/:id/xp", h.SetXP())
}