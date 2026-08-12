package handlers

import (

	"net/http"
	"github.com/gin-gonic/gin"
	"backend/internal/models/games"
	"backend/internal/repositories"
)

type LobbyHandler struct {
	gameRepo *repositories.GameRepository
	lobbyRepo *repositories.LobbyRepository
}

func NewLobbyHandler(gameRepo *repositories.GameRepository,
	lobbyRepo *repositories.lobbyRepository) *GameHandler {
	return &GameHandler{gameRepo: gameRepo, lobbyRepo: lobbyRepo}
}

func (h *LobbyHandler) CreateLobby() gin.HandlerFunc {
	return func(c *gin.Context) {
		gameObj := &games.Game{}
		game, err := h.gameRepo.CreateGame(c.Request.Context(), gameObj)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "the creation of the game failed, try again"})
			return ;
		}		
		lobbyObj := &games.Lobby{UserId: user_id,
								GameId: game.ID}
		lobby, err := h.lobbyRepo.CreateLobby(c.Request.Context(), lobbyObj)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "the creation of the lobby failed, try again"})
			return ;
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success"})
	}
}