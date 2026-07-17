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
		public.POST("/user/create", user.CreateUser(db))
		public.POST("/auth/login", user.LoginUser(db))
		// public.POST("/auth/logout", user.Logout(db))
	}

	// it's private because only a registred user can access it
	private := r.Group("/")
	private.Use(auth.AuthMiddleWare())
	{
		private.GET("/user/get/:id", user.GetUser(db))
		private.GET("/user/get/name/:id", user.GetName(db))
		private.GET("/user/get/login/:id", user.GetLogin(db))
		private.GET("/user/get/email/:id", user.GetEmail(db))
		private.GET("/user/get/oauth/:id", user.GetOAuth(db))
		// private.GET("/user/get/avatar/:id", user.GetAvatar(db))
		

		// PUT = modify a lot of informations
		// PATCH = modify a specific field
		private.PATCH("/user/name/:id", user.SetName(db))
		private.PATCH("/user/login/:id", user.SetLogin(db))
		private.PATCH("/user/avatar/:id", user.SetAvatar(db))
	}
	return r, nil
}