package router

import (
    // EXTERNAL PACKAGE
    "os"
	"github.com/uptrace/bun"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    
    // OUR PACKAGE
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
    r := gin.Default()
    env := os.Getenv("ENV")
    if len(env) != 0 && env == "prod" {
        r.Use(cors.New(cors.Config{
            AllowOrigins:     []string{"http://react:80/"},
            AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
            ExposeHeaders:    []string{"Content-Length"},
            AllowCredentials: true,
        }))
    } else {    
        r.Use(cors.Default())
    }
	public := r.Group("/")
    RegisterUserPublicRoutes(public, db)
    RegisterAuthPublicRoutes(public, db)

	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
    RegisterUserPrivateRoutes(private, db)

	return r, nil
}
