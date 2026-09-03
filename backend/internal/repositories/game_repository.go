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

func (r *GameRepository) Game(ctx context.Context, game_id int64) (*games.Game, error){
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

func (r *GameRepository) Title(ctx context.Context, game_id int64) (string, error){
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

func (r *GameRepository) Mode(ctx context.Context, game_id int64) (string, error){
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

func (r *GameRepository) Difficulty(ctx context.Context, game_id int64) (string, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("difficulty").
		Where("id = ?", game_id).
		Scan(ctx)
	if err != nil {
		return "", err
	}
	return game.Difficulty, nil
}

func (r *GameRepository) Timer(ctx context.Context, game_id int64) (int64, error){
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

func (r *GameRepository) RoundTotal(ctx context.Context, game_id int64) (int64, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("round_total").
		Where("id = ?", game_id).
		Scan(ctx)
	if err != nil {
		return -1, err
	}
	return game.RoundTotal, nil
}

func (r *GameRepository) Finished(ctx context.Context, game_id int64) (bool, error){
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

func (r *GameRepository) PlayerNb(ctx context.Context, game_id int64) (int64, error){
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

func (r *GameRepository) CreatedAt(ctx context.Context, game_id int64) (time.Time, error){
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

func (r *GameRepository) SetMode(ctx context.Context, game_id int64, 
	mode string) error {
	game := &games.Game{Mode: mode}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("mode").
		Where("id = ?", game_id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) SetDifficulty(ctx context.Context, game_id int64, 
	difficulty string) error {
	game := &games.Game{Difficulty: difficulty}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("difficulty").
		Where("id = ?", game_id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) SetTimer(ctx context.Context, game_id int64, 
	timer int64) error {
	game := &games.Game{Timer: timer}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("timer").
		Where("id = ?", game_id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) SetRoundTotal(ctx context.Context, game_id int64, 
	round_total int64) error {
	game := &games.Game{RoundTotal: round_total}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("round_total").
		Where("id = ?", game_id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) SetFinished(ctx context.Context, game_id int64, 
	finished bool) error {
	game := &games.Game{Finished: finished}
	_, err := r.db.NewUpdate().
		Model(game).
		Column("finished").
		Where("id = ?", game_id).
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

// GetRandomLocations returns a list of unique random location IDs based on the total number of rounds
func (r *LocationRepository) GetRandomLocations(ctx context.Context, round_total int64) ([]int64, error) {
	locations := make([]int64, 0)
	err := r.db.NewSelect().
		Model(&locations).
		TableExpr("locations").
		OrderExpr("RANDOM()").
		Limit(int(round_total)). // limits to the number of rounds
		Group("location_id"). // ensures unique locations
		Column("location_id"). // select only the location id column
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return locations, nil
}
