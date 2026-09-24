package router

import (
	"backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func registerRoundRoutes(r *gin.Engine, h *handlers.RoundHandler) {
	// Public user
	public := r.Group("/")
	{
		public.POST("/rounds/", h.NextRound) //create the new round using its round nb (user_id and game_id should be in the context/token)


		// public.GET("/rounds/:id", h.GetRoundByID) //get a round by id

		// public.PATCH("/rounds/:id", h.UpdateRoundByID) //update a round  by id (game_id, round_nb, location_id, start_time, end_time, duration, is_active)

		// public.DELETE("/rounds/:id", h.DeleteRoundByID) //delete a round by id
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