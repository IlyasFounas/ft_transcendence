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
	public := r.Group("/")
	{
		public.POST("/user/create", h.CreateUser())
		public.POST("/auth/login", h.LoginUser())
		// BUGFIX: ces deux routes appelaient h.ActivationToken() / h.RefreshToken(),
		// des méthodes qui n'existent pas sur UserHandler — elles vivent sur TokenHandler.
		public.GET("/user/activation", ht.ActivationToken())
		public.GET("/auth/refresh", ht.RefreshToken())
	}
	// it's private because only a registred user can access it
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.GET("/user/get", h.GetUser())
		private.GET("/user/get/name", h.GetName())
		private.GET("/user/get/login", h.GetLogin())
		private.GET("/user/get/email", h.GetEmail())
		private.GET("/user/get/oauth", h.GetOAuth())
		private.GET("/user/:username", h.UserByName())
		// PUT = modify a lot of informations
		// PATCH = modify a specific field
		private.PATCH("/user/name", h.SetName())
		private.PATCH("/user/login", h.SetLogin())
		private.PATCH("/user/avatar", h.SetAvatar())
	}
}
