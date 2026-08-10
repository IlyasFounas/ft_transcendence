package repositories

import (
	"time"
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

func (r *GameRepository) GetGame(ctx context.Context, game_id int64) (*games.Game, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Where("id = ?", game_id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (r *GameRepository) GetTitle(ctx context.Context, game_id int64) (string, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("title").
		Where("id = ?", game_id).
		Scan(ctx)
	if err != nil {
		return "", err
	}
	return game.Title, nil
}

func (r *GameRepository) GetMode(ctx context.Context, game_id int64) (string, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("mode").
		Where("id = ?", game_id).
		Scan(ctx)
	if err != nil {
		return "", err
	}
	return game.Mode, nil
}

func (r *GameRepository) GetTimer(ctx context.Context, game_id int64) (int64, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("timer").
		Where("id = ?", game_id).
		Scan(ctx)
	if err != nil {
		return -1, err
	}
	return game.Timer, nil
}

func (r *GameRepository) GetRoundNb(ctx context.Context, game_id int64) (int64, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("round_nb").
		Where("id = ?", game_id).
		Scan(ctx)
	if err != nil {
		return -1, err
	}
	return game.RoundNb, nil
}

func (r *GameRepository) GetFinished(ctx context.Context, game_id int64) (bool, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("finished").
		Where("id = ?", game_id).
		Scan(ctx)
	if err != nil {
		return false, err
	}
	return game.Finished, nil
}

func (r *GameRepository) GetPlayerNb(ctx context.Context, game_id int64) (int64, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("player_nb").
		Where("id = ?", game_id).
		Scan(ctx)
	if err != nil {
		return -1, err
	}
	return game.PlayerNb, nil
}

func (r *GameRepository) GetCreatedAt(ctx context.Context, game_id int64) (time.Time, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("created_at").
		Where("id = ?", game_id).
		Scan(ctx)
	if err != nil {
		return time.Now(), err
	}
	return game.CreatedAt, nil
}

func (r *GameRepository) SetTitle(ctx context.Context, game_id int64, 
	title string) error {
	game := &games.Game{Title: title}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("title").
		Where("id = ?", game_id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}