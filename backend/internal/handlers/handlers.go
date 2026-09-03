package handlers

import (
	// OUR PACKAGE
	"backend/internal/repositories"
)

type Handlers struct {
	Location *LocationHandler
	User     *UserHandler
	Token    *TokenHandler
	Game	 *GameHandler
	Lobby	 *LobbyHandler
	Round	 *RoundHandler
	Guess	 *GuessHandler
}

// NewHandlers instanciates and assembles all the handlers
func NewHandlers(repo *repositories.Repositories) *Handlers {
	tokenHandler := NewTokenHandler(repo.User, repo.AToken, repo.RToken)
	return &Handlers{
		Location: NewLocationHandler(repo.Location),
		User:     NewUserHandler(repo.User, repo.AToken, tokenHandler),
		Token:    tokenHandler,
		Game:     NewGameHandler(repo.Game, repo.Lobby, repo.Location),
		Lobby:    NewLobbyHandler(repo.Game, repo.Lobby, repo.User),
		Round:    NewRoundHandler(repo.Round),
		Guess:    NewGuessHandler(repo.Guess),
	}
}
