package router

import (
	// "log" // TEMPORARY TEST SHUTDOWN
	// "net/http" // TEMPORARY TEST SHUTDOWN
	// "time" // TEMPORARY TEST SHUTDOWN

	"backend/internal/handlers"
	"backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func registerLocationRoutes(r *gin.Engine, h *handlers.LocationHandler) {
	// Public user
	// public := r.Group("/")
	// {
	// // TEMPORARY TEST ROUTE FOR GRACEFUL SHUTDOWN
	// public.GET("locations/slow", func(c *gin.Context) {
	// 	log.Println("[Request] Started processing slow endpoint...")
	// 	time.Sleep(5 * time.Second) // Simulates heavy work
	// 	c.String(http.StatusOK, "Finished!")
	// 	log.Println("[Request] Completed slow endpoint.")
	// })
	// }

	admin := r.Group("/")
	admin.Use(middleware.AuthMiddleWare(), middleware.AdminMiddleWare())
	{
		admin.POST("/admin/locations", h.AddLocations)
		admin.GET("/admin/locations/:location_id", h.GetLocationByID)
		admin.GET("/admin/locations/mapillary/:mid", h.GetLocationByMID)
		admin.GET("/admin/locations/random", h.GetRandomLocation)
		admin.GET("/admin/locations/random/country/:cc", h.GetRandomLocationByCC)
		admin.GET("/admin/locations/:location_id/coordinates", h.GetCoordinatesByID)
		admin.GET("/admin/locations/:location_id/cc", h.GetCountryCodeByID)
		admin.PATCH("/admin/locations/:location_id", h.UpdateLocationByID)
		admin.PATCH("/admin/locations/mapillary/:mid", h.UpdateLocationByMID)
		admin.DELETE("/admin/locations/:location_id", h.DeleteLocation)
		admin.DELETE("/admin/locations/mapillary/:mid", h.DeleteLocationByMID)
	}
}
