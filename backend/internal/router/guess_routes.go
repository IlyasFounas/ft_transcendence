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
		admin.GET("/admin/guesses/:id", h.GetGuessByID) //get a guess by id
		admin.GET("/admin/guesses/:id/coordinates", h.GetGuessCoordinates) //get coordinates (lat + long) for a guess
		admin.GET("/admin/guesses/:id/latitude", h.GetGuessLatitude) //get latitude for a guess
		admin.GET("/admin/guesses/:id/longitude", h.GetGuessLongitude) //get longitude for a guess
		admin.GET("/admin/guesses/:id/country", h.GetGuessCountryCode) //get country_code for a guess
		admin.GET("/admin/guesses/:id/distance", h.GetGuessDistance) //get distance for a guess
		admin.GET("/admin/guesses/:id/score", h.GetGuessScore) //get score for a guess
		admin.GET("/admin/guesses/score/:round_id/:user_id", h.GetGuessScoreByRoundAndUser) //get score for a guess by round and user
		admin.GET("/admin/guesses/results/:game_id", h.GetGameResults) //get results for a game (map of map with user_id, total_score)
		admin.GET("/guesses/total_score/:game_id/:user_id", h.GetTotalScoreByGameAndUser) //get total score for a user in a game
		admin.DELETE("/admin/guesses/:id", h.DeleteGuessByID)
	}
}