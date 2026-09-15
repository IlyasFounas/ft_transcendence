package ws_handlers

import (
	"backend/internal/repositories"
	"github.com/sams96/rgeo"
)

type WsHandlers struct {
	Lobby *LobbyWsHandler
	Guess *GuessWsHandler
}

func NewWsHandlers(repo *repositories.Repositories, reverser *rgeo.Rgeo) *WsHandlers {
	return &WsHandlers{
		Lobby: NewLobbyWsHandler(repo.Lobby, repo.User),
		Guess: NewGuessWsHandler(repo.Guess, repo.Round, repo.Lobby, repo.Game, repo.User, reverser),
	}
}