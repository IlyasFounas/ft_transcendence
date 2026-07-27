package router

import (
	"github.com/gin-gonic/gin"
	"backend/internal/handler"
)

func registerLocationRoutes(r *gin.Engine, h *handler.LocationHandler) {
	// Public user
	public := r.Group("/")
	{
		public.POST("/locations", h.AddLocations) // TODO: put in admin group when admin exists
		public.GET("/locations/:id", h.GetLocationByID)
		public.GET("/locations/mapillary/:mid", h.GetLocationByMID)
		public.GET("/locations/random", h.GetRandomLocation)
		public.GET("/locations/random/country/:cc", h.GetRandomLocationByCC)
		public.PATCH("/locations/:id", h.UpdateLocationByID) // TODO: admin
		public.PATCH("/locations/mapillary/:mid", h.UpdateLocationByMID) // TODO: admin
		public.DELETE("/locations/:id", h.DeleteLocation) // TODO: admin
		public.DELETE("/locations/mapillary/:mid", h.DeleteLocationByMID) // TODO: admin
	}

// 	// Authentificated user
// 	// private := r.Group("/")
// 	// private.Use(auth.AuthRequired())
// 	// {

// 	// }

// 	// Admin user
// 	// admin := r.Group("/")
// 	// admin.Use(auth.AuthRequired())
// 	// {

// 	// }
// 	return r, nil
}