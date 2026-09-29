package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
)

type Repositories struct {
	Location	*LocationRepository
	User		*UserRepository
	AToken		*ActivationTokenRepository
	RToken		*RefreshTokenRepository
	Game		*GameRepository
	Lobby		*LobbyRepository
	Round		*RoundRepository
	Guess		*GuessRepository
	Friend		*FriendRepository
	Chat		*ChatRepository
	Statistics	*StatisticsRepository
}

// NewRepositories instanciates and assembles all the repositories
func NewRepositories(db *bun.DB) (*Repositories, error) {
	locationRepo, err := NewLocationRepository(db)
	if err != nil {
		return nil, err
	}

	userRepo := NewUserRepository(db)

	return &Repositories{
		Location:	locationRepo,
		User:		userRepo,
		AToken:		NewActivationTokenRepository(db),
		RToken:		NewRefreshTokenRepository(db),
		Game:		NewGameRepository(db),
		Lobby:		NewLobbyRepository(db),
		Round:		NewRoundRepository(db, locationRepo),
		Guess:		NewGuessRepository(db),
		Friend:		NewFriendRepository(db),
		Chat:		NewChatRepository(db),
		Statistics:	NewStatisticsRepository(db),
	}, nil
}

// Update for each user that stayed until the end of this game his xp, then his rank from this xp, then his statistics
func (r *Repositories) UpdatePlayers(ctx context.Context, userIds []int64) error {
	var errs []error
	if err := r.User.UpdateXP(ctx, userIds); err != nil {
		errs = append(errs, fmt.Errorf("xp: %w", err))
	}
	if err := r.User.UpdateRank(ctx, userIds); err != nil {
		errs = append(errs, fmt.Errorf("rank: %w", err))
	}
	if err := r.Statistics.UpdateStats(ctx, userIds); err != nil {
		errs = append(errs, fmt.Errorf("statistics: %w", err))
	}
	return errors.Join(errs...)
}

// Runs UpdatePlayers for every user
func (r *Repositories) UpdateAllPlayers(ctx context.Context) error {
	userIds, err := r.User.AllIds(ctx)
	if err != nil {
		return err
	}
	return r.UpdatePlayers(ctx, userIds)
}