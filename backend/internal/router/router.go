package router

import (
	// EXTERNAL PACKAGE
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// OUR PACKAGE
	"backend/internal/handlers"
)

func NewRouter(h *handlers.Handlers) *gin.Engine {
	r := gin.Default()    //when you first create a var :=
	r.Use(cors.Default()) // TODO: global middleware now but accept only localhost:3000 request

	// Splitting routes
	// registerUserRoutes(r, h.User)
	registerLocationRoutes(r, h.Location)

	return r
}
