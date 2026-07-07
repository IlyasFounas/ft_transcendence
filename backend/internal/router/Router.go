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

func InitRouter(db *bun.DB) (*gin.Engine, error) {
    r := gin.Default() //when you first create a var :=
    r.Use(cors.Default()) // TODO accept only localhost:3000 request

	public := r.Group("/")
	{
		public.GET("/", user.HomePage)
		public.POST("/user/login", user.LoginUser)
	}

	// it's private because only a registred user can access it
	private := r.Group("/")
	private.Use(auth.AuthRequired())
	{
		private.GET("/user/:id", user.GetUser)
	}
	return r, nil
}