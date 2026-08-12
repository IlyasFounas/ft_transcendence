package repositories

import (
	"github.com/uptrace/bun"
)

type Repositories struct {
	Location *LocationRepository
	User     *UserRepository
	AToken 	 *ActivationTokenRepository
	RToken   *RefreshTokenRepository
	Game     *GameRepository
	Lobby    *LobbyRepository
}

// NewRepositories instanciates and assembles all the repositories
func NewRepositories(db *bun.DB) (*Repositories, error) {
	locationRepo, err := NewLocationRepository(db)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		Location: locationRepo,
		User:     NewUserRepository(db),
		AToken:   NewActivationTokenRepository(db),
		RToken:   NewRefreshTokenRepository(db),
		Game:     NewGameRepository(db),
		Lobby:    NewLobbyRepository(db),
	}, nil
}
