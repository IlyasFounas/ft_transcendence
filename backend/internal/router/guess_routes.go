package router

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func registerGuessRoutes(r *gin.Engine, h *handlers.GuessHandler) {
	// Public user
	public := r.Group("/")
	{
		public.GET("/guesses/:id", h.GetGuessByID) //get a guess by id
		public.GET("/guesses/:id/coordinates", h.GetGuessCoordinates) //get coordinates (lat + long) for a guess
		public.GET("/guesses/:id/latitude", h.GetGuessLatitude) //get latitude for a guess
		public.GET("/guesses/:id/longitude", h.GetGuessLongitude) //get longitude for a guess
		public.GET("/guesses/:id/country", h.GetGuessCountryCode) //get country_code for a guess
		public.GET("/guesses/:id/distance", h.GetGuessDistance) //get distance for a guess
		public.GET("/guesses/:id/score", h.GetGuessScore) //get score for a guess
		public.GET("/guesses/leaderboard/:game_id", h.GetLeaderboardByGame) //get leaderboard for a game (map of map with user_id, total_score)

		// public.PATCH("/guesses/:id", h.UpdateGuessByID) //update a guess  by id (latitude, longitude, round_id, user_id, game_id)

		public.DELETE("/guesses/:id", h.DeleteGuessByID) //delete a guess by id
		// public.DELETE("/guesses/round/:round_id", h.DeleteGuessesByRoundID) //delete all guesses for a round
		// public.DELETE("/guesses/user/:user_id", h.DeleteGuessesByUserID) //delete all guesses for a user
		// public.DELETE("/guesses/game/:game_id", h.DeleteGuessesByGameID) //delete all guesses of all users for a game
	}

	private := r.Group("/")
	private.Use(middleware.AuthMiddleWare())
	{
		private.GET("/guesses/xp/:user_id", h.GetUserXP) //get total of all games of a user = total XP of a user
		private.GET("/guesses/score/:round_id/:user_id", h.GetGuessScoreByRoundAndUser) //get score for a guess by round and user
		private.GET("/guesses/total_score/:game_id/:user_id", h.GetTotalScoreByGameAndUser) //get total score for a user in a game
	}

	// 	// Admin user
	// 	// admin := r.Group("/")
	// 	// admin.Use(auth.AuthRequired())
	// 	// {

	// 	// }
}