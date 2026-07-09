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
 */

func InitRouter(db *bun.DB) (*gin.Engine, error) {
    r := gin.Default() //when you first create a var :=
    r.Use(cors.Default()) // TODO accept only localhost:3000 request

	public := r.Group("/")
	{
		public.GET("/", user.HomePage)
		public.POST("/user/login", user.LoginUser(db))
		public.POST("/user/create", user.CreateUser(db))
	}

	// it's private because only a registred user can access it
	private := r.Group("/")
	private.Use(auth.AuthRequired())
	{
		private.GET("/profile/:id", user.GetUser)
	}
	return r, nil
}