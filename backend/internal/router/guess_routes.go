package router

import (
	"backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func registerGuessRoutes(r *gin.Engine, h *handlers.GuessHandler) {
	// Public user
	public := r.Group("/")
	{
		public.POST("/guesses", h.AddGuess) //add a guess (latitude, longitude, round_id?, user_id?, game_id?)
		// calculate cc based on lat + long (using reverse geocoding API)
		// calculate distance based on lat + long compared to actual location
		// calculate score based on distance and correctness of country code (add basic scoring + multipliers for correcteness of cc and streaks, and divider for time taken)
		// return guess with all calculated fields (cc, distance, score)

		public.GET("/guesses/:id", h.GetGuessByID) //get a guess by id
		public.GET("/guesses/:id/coordinates", h.GetGuessCoordinates) //get coordinates (lat + long) for a guess
		public.GET("/guesses/:id/latitude", h.GetGuessLatitude) //get latitude for a guess
		public.GET("/guesses/:id/longitude", h.GetGuessLongitude) //get longitude for a guess
		public.GET("/guesses/:id/country", h.GetGuessCountryCode) //get country_code for a guess
		public.GET("/guesses/:id/distance", h.GetGuessDistance) //get distance for a guess
		public.GET("/guesses/:id/score", h.GetGuessScore) //get score for a guess
		public.GET("/guesses/score/:round_id/:user_id", h.GetGuessScoreByRoundAndUser) //get score for a guess by round and user // ou /guesses/score/:round_id/:user_id
		public.GET("/guesses/total_score/:game_id/:user_id", h.GetTotalScoreByGameAndUser) //get total score for a user in a game
		public.GET("/guesses/xp/:user_id", h.GetUserXP) //get total of all games of a user = total XP of a user
		public.GET("/guesses/leaderboard/:game_id", h.GetLeaderboardByGame) //get leaderboard for a game (map of map with user_id, total_score)
		
		// public.PATCH("/guesses/:id", h.UpdateGuessByID) //update a guess  by id (latitude, longitude, round_id, user_id, game_id)

		public.DELETE("/guesses/:id", h.DeleteGuessByID) //delete a guess by id
		// public.DELETE("/guesses/round/:round_id", h.DeleteGuessesByRoundID) //delete all guesses for a round
		// public.DELETE("/guesses/user/:user_id", h.DeleteGuessesByUserID) //delete all guesses for a user
		// public.DELETE("/guesses/game/:game_id", h.DeleteGuessesByGameID) //delete all guesses of all users for a game
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
}
