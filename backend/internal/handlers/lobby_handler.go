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
	lobbyRepo *repositories.LobbyRepository) *LobbyHandler {
	return &LobbyHandler{gameRepo: gameRepo, lobbyRepo: lobbyRepo}
}

// id, err := strconv.ParseInt(c.Param("Join"), 10, 64)
func (h *LobbyHandler) JoinLobby() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		joinCode := c.PostForm("join_code")
		gameId, err := h.lobbyRepo.GameIdLobby(c.Request.Context(), joinCode)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return
		}
		playerNb, err := h.lobbyRepo.LobbyPlayerNb(c.Request.Context(), gameId)
		if err != nil{
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return
		}
		gamePlayerNb, err := h.gameRepo.PlayerNb(c.Request.Context(), gameId)
		if playerNb + 1 > int(gamePlayerNb) {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The lobby is full"})
			return
		}
		lobbyObj := &games.Lobby{UserId: userId.(int64),
								GameId: gameId}
		err = h.lobbyRepo.JoinLobby(c.Request.Context(), lobbyObj)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success"})
	}
}
