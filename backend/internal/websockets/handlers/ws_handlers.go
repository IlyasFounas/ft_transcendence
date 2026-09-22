package ws_handlers

import (
	"backend/internal/repositories"
	"github.com/sams96/rgeo"
)

type WsHandlers struct {
	Lobby	*LobbyWsHandler
	Round	*RoundWsHandler
	Guess	*GuessWsHandler
	Chat	*ChatHandler
}

func NewWsHandlers(repo *repositories.Repositories, reverser *rgeo.Rgeo) *WsHandlers {
	return &WsHandlers{
		Lobby: NewLobbyWsHandler(repo.Lobby, repo.Game, repo.User),
		Round: NewRoundWsHandler(repo.Round),
		Guess: NewGuessWsHandler(repo.Guess, repo.Round, repo.Lobby, repo.Game, repo.User, reverser),
		Chat: NewChatHandler(repo.Chat, repo.Friend, repo.User),
	}
}