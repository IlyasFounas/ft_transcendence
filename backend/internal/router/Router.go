package router

import (
    // EXTERNAL PACKAGE
	"github.com/uptrace/bun"
    "github.com/gin-contrib/cors"
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

func InitRouter(db *bun.DB) (*gin.Engine, error) {
    r := gin.Default() //when you first create a var :=
    r.Use(cors.Default()) // TODO accept only localhost:3000 request

	public := r.Group("/")
	{
		public.POST("/user/create", handler.CreateUser(db))
		public.POST("/auth/login", handler.LoginUser(db))
		public.GET("/user/activation", handler.UserActivation(db))

		// public.POST("/auth/logout", user.Logout(db))
	}

	// it's private because only a registred user can access it
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.GET("/user/get", handler.GetUser(db))
		private.GET("/user/get/name", handler.GetName(db))
		private.GET("/user/get/login", handler.GetLogin(db))
		private.GET("/user/get/email", handler.GetEmail(db))
		private.GET("/user/get/oauth", handler.GetOAuth(db))
		// private.GET("/user/get/avatar/:id", user.GetAvatar(db))
		

		// PUT = modify a lot of informations
		// PATCH = modify a specific field
		private.PATCH("/user/name", handler.SetName(db))
		private.PATCH("/user/login", handler.SetLogin(db))
		private.PATCH("/user/avatar", handler.SetAvatar(db))
	}
	return r, nil
}