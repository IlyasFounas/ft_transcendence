package router

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func registerGuessRoutes(r *gin.Engine, h *handlers.GuessHandler) {

	admin := r.Group("/")
	admin.Use(middleware.AuthMiddleWare(), middleware.AdminMiddleWare())
	{
		admin.GET("/admin/guesses/:id", h.GetGuessByID)
		admin.GET("/admin/guesses/:id/coordinates", h.GetGuessCoordinates)
		admin.GET("/admin/guesses/:id/latitude", h.GetGuessLatitude)
		admin.GET("/admin/guesses/:id/longitude", h.GetGuessLongitude)
		admin.GET("/admin/guesses/:id/country", h.GetGuessCountryCode)
		admin.GET("/admin/guesses/:id/distance", h.GetGuessDistance)
		admin.GET("/admin/guesses/:id/score", h.GetGuessScore)
		admin.GET("/admin/guesses/score/:round_id/:user_id", h.GetGuessScoreByRoundAndUser)
		admin.GET("/admin/guesses/results/:game_id", h.GetGameResults)
		admin.GET("/admin/guesses/total_score/:game_id/:user_id", h.GetTotalScoreByGameAndUser)
		admin.DELETE("/admin/guesses/:id", h.DeleteGuessByID)
	}
}
