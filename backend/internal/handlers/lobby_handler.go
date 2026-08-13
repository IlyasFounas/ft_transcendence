package handlers

import (

	// "net/http"
	// "github.com/gin-gonic/gin"
	// "backend/internal/models/games"
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

