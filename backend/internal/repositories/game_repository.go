package repositories

import (
	"context"
	"time"

	// "errors"
	"backend/internal/models/games"

	"github.com/uptrace/bun"
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

// Getters
func (r *GameRepository) Game(ctx context.Context, gameId int64) (*games.Game, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Where("id = ?", gameId).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (r *GameRepository) Title(ctx context.Context, gameId int64) (string, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("title").
		Where("id = ?", gameId).
		Scan(ctx)
	if err != nil {
		return "", err
	}
	return game.Title, nil
}

func (r *GameRepository) Mode(ctx context.Context, gameId int64) (string, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("mode").
		Where("id = ?", gameId).
		Scan(ctx)
	if err != nil {
		return "", err
	}
	return game.Mode, nil
}

func (r *GameRepository) Difficulty(ctx context.Context, gameId int64) (string, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("difficulty").
		Where("id = ?", gameId).
		Scan(ctx)
	if err != nil {
		return "", err
	}
	return game.Difficulty, nil
}

func (r *GameRepository) Timer(ctx context.Context, gameId int64) (int64, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("timer").
		Where("id = ?", gameId).
		Scan(ctx)
	if err != nil {
		return -1, err
	}
	return game.Timer, nil
}

func (r *GameRepository) RoundNb(ctx context.Context, gameId int64) (int64, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("round_nb").
		Where("id = ?", gameId).
		Scan(ctx)
	if err != nil {
		return -1, err
	}
	return game.RoundNb, nil
}

func (r *GameRepository) Finished(ctx context.Context, gameId int64) (*time.Time, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("finished").
		Where("id = ?", gameId).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return game.FinishedAt, nil
}

func (r *GameRepository) PlayerNb(ctx context.Context, gameId int64) (int64, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("player_nb").
		Where("id = ?", gameId).
		Scan(ctx)
	if err != nil {
		return -1, err
	}
	return game.PlayerNb, nil
}

func (r *GameRepository) CreatedAt(ctx context.Context, gameId int64) (time.Time, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("created_at").
		Where("id = ?", gameId).
		Scan(ctx)
	if err != nil {
		return time.Now(), err
	}
	return game.CreatedAt, nil
}

// Setters
func (r *GameRepository) SetTitle(ctx context.Context, gameId int64, 
	title string) error {
	game := &games.Game{Title: title}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("title").
		Where("id = ?", gameId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) SetMode(ctx context.Context, gameId int64, 
	mode string) error {
	game := &games.Game{Mode: mode}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("mode").
		Where("id = ?", gameId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) SetDifficulty(ctx context.Context, gameId int64, 
	difficulty string) error {
	game := &games.Game{Difficulty: difficulty}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("difficulty").
		Where("id = ?", gameId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) SetTimer(ctx context.Context, gameId int64, 
	timer int64) error {
	game := &games.Game{Timer: timer}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("timer").
		Where("id = ?", gameId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) SetRoundNb(ctx context.Context, gameId int64, 
	round_nb int64) error {
	game := &games.Game{RoundNb: round_nb}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("round_nb").
		Where("id = ?", gameId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) SetFinished(ctx context.Context,
	gameId int64) error {
	finished := time.Now()
	game := &games.Game{FinishedAt: &finished}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("finished_at").
		Where("id = ?", gameId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) SetPlayerNb(ctx context.Context, game_id int64, 
	player_nb int64) error {
	game := &games.Game{PlayerNb: player_nb}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("player_nb").
		Where("id = ?", game_id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}