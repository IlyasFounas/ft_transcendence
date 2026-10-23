package repositories

import (
	"github.com/uptrace/bun"
)

type Repositories struct {
	Location   *LocationRepository
	User       *UserRepository
	AToken     *ActivationTokenRepository
	RToken     *RefreshTokenRepository
	Game       *GameRepository
	Lobby      *LobbyRepository
	Round      *RoundRepository
	Guess      *GuessRepository
	Friend     *FriendRepository
	Chat       *ChatRepository
	Statistics *StatisticsRepository
}

// NewRepositories instanciates and assembles all the repositories
func NewRepositories(db *bun.DB) (*Repositories, error) {
	locationRepo, err := NewLocationRepository(db)
	if err != nil {
		return nil, err
	}

	userRepo := NewUserRepository(db)

	return &Repositories{
		Location:   locationRepo,
		User:       userRepo,
		AToken:     NewActivationTokenRepository(db),
		RToken:     NewRefreshTokenRepository(db),
		Game:       NewGameRepository(db),
		Lobby:      NewLobbyRepository(db),
		Round:      NewRoundRepository(db, locationRepo),
		Guess:      NewGuessRepository(db),
		Friend:     NewFriendRepository(db),
		Chat:       NewChatRepository(db),
		Statistics: NewStatisticsRepository(db),
	}, nil
}
