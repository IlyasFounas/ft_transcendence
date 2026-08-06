package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"backend/internal/models/games"
	"backend/internal/repositories"
)

type GameHandler struct {
	gameRepo *repositories.GameRepository
} 

func NewGameHandler(gameRepo *repositories.GameRepository) *GameHandler {
	return &GameHandler{gameRepo: gameRepo}
}

func (h *GameHandler) CreateGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		game = games.Game{}
		game, err := gameRepo.CreateGame(c.Request.Context(), )
		if err != nil {
			c.JSON(http.StatusOK,
				gin.H{"message": "the creation of the game failed, try again"})
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "the game is launching", "game_id": game.ID})
	}
}
