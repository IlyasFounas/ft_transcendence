package handlers

import (
	"backend/internal/repositories"
	"github.com/sams96/rgeo"
	"os"
	"backend/internal/repositories"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Handlers struct {
	Location	*LocationHandler
	User		*UserHandler
	Token		*TokenHandler
	Game		*GameHandler
	Lobby		*LobbyHandler
	Round		*RoundHandler
	Guess		*GuessHandler
}

func googleConfigInit() *oauth2.Config {
	googleConf := &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	Scopes:       []string{"openid", "email", "profile"},
	Endpoint:     google.Endpoint,}
	return googleConf
}

// NewHandlers instanciates and assembles all the handlers
func NewHandlers(repo *repositories.Repositories, reverser *rgeo.Rgeo) *Handlers {
	tokenHandler := NewTokenHandler(repo.User, repo.AToken, repo.RToken)
	googleConfig := googleConfigInit()
	return &Handlers{
		Location:	NewLocationHandler(repo.Location),
		User:	 	NewUserHandler(repo.User, repo.AToken, tokenHandler, googleConfig),
		Token:		tokenHandler,
		Game:		NewGameHandler(repo.Game, repo.Lobby, repo.Location),
		Lobby:		NewLobbyHandler(repo.Game, repo.Lobby, repo.User),
		Round:		NewRoundHandler(repo.Round),
		Guess:		NewGuessHandler(repo.Guess, repo.Round, repo.Lobby, repo.Game, reverser),
	}
}
