package router

import (
    // EXTERNAL PACKAGE
	"github.com/uptrace/bun"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    
    // OUR PACKAGE
    "backend/internal/handler"
	"backend/internal/repository"
)

func InitRouter(db *bun.DB) *gin.Engine {
	r := gin.Default() //when you first create a var :=
	r.Use(cors.Default()) // TODO: global middleware now but accept only localhost:3000 request

	// Init dependencies
	// userRepo := repository.NewUserRepository(db)
	// userHandler := handler.NewUserHandler(userRepo)

	locationRepo := repository.NewLocationRepository(db)
	locationHandler := handler.NewLocationHandler(locationRepo)

	// Splitting routes
	// registerUserRoutes(r, userHandler)
	registerLocationRoutes(r, locationHandler)

	return r // no need for nil error return since nothing done in this function can fail (only not compile)
}