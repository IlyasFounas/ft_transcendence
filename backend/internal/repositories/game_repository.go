package repositories

import (

	"context"
	// "errors"
	"github.com/uptrace/bun"
	"backend/internal/models/games"
)

type GameRepository struct {
	db *bun.DB
}

func NewGameRepository(db *bun.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) CreateGame(ctx context.Context, game *games.Game) (*games.Game, error){
	_, err := r.db.NewInsert().
		Model(game).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return game, nil
}