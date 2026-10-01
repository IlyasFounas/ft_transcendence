package router

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func registerStatisticsRoutes(r *gin.Engine, h *handlers.StatisticsHandler) {
	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		// statistics of the connected user
		private.GET("/user/statistics", h.Statistics()) // all of them
		private.GET("/user/statistics/finishedgames", h.FinishedGames())
		private.GET("/user/statistics/accuracy", h.Accuracy())
		private.GET("/user/statistics/bestscore", h.BestScore())
		private.GET("/user/statistics/closestguess", h.ClosestGuess())
	}

	admin := r.Group("/")
	admin.Use(middleware.AuthMiddleWare(), middleware.AdminMiddleWare())
	{
		// admin.PATCH("/admin/user/:id/statistics/finishedgames", h.SetFinishedGames())
		// admin.PATCH("/admin/user/:id/statistics/accuracy", h.SetAccuracy())
		// admin.PATCH("/admin/user/:id/statistics/bestscore", h.SetBestScore())
		// admin.PATCH("/admin/user/:id/statistics/closestguess", h.SetClosestGuess())
	}
}