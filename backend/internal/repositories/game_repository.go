package repositories

import (
	"backend/internal/models/games"
	"context"
	"errors"
	"github.com/uptrace/bun"
	"time"
)

type GameRepository struct {
	db       *bun.DB
	userRepo *UserRepository // to add the xp of the players when the game is finished
}

func NewGameRepository(db *bun.DB, userRepo *UserRepository) *GameRepository {
	return &GameRepository{db: db, userRepo: userRepo}
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

func (r *GameRepository) RoundTotal(ctx context.Context, gameId int64) (int64, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("round_total").
		Where("id = ?", gameId).
		Scan(ctx)
	if err != nil {
		return -1, err
	}
	return game.RoundTotal, nil
}

func (r *GameRepository) Finished(ctx context.Context, gameId int64) (bool, error){
	return r.db.NewSelect().
		Model((*games.Game)(nil)).
		Column("finished_at").
		Where("id = ?", gameId).
		Exists(ctx)
}

func (r *GameRepository) FinishedAt(ctx context.Context, gameId int64) (*time.Time, error){
	game := new(games.Game)
	err := r.db.NewSelect().
		Model(game).
		Column("finished_at").
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

// First call sets finished_at and, in the same transaction,
// adds to each player the base_score he earned in this game (users.xp)
// The next calls do nothing, so the xp of a game is never counted twice
// (the hub, the HTTP route and a reconnection can all call it).
func (r *GameRepository) SetFinished(ctx context.Context,
	gameId int64) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model((*games.Game)(nil)).
			Set("finished_at = ?", time.Now()).
			Where("id = ?", gameId).
			Where("finished_at IS NULL").
			Exec(ctx)
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			// already finished: the xp of this game has already been added
			return nil
		}
		return r.userRepo.AddGameXP(ctx, tx, gameId)
	})
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

// Check if round_total is smaller than the number of available locations (else error)
func (r *LocationRepository) CheckEnoughLocations(ctx context.Context, round_total int64) (error) {
	var count int
	count, err := r.db.NewSelect().
		TableExpr("locations").
		Count(ctx)
	if err != nil {
		return err
	}
	if int(round_total) > count {
		return errors.New("not enough unique locations available")
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
		Limit(int(round_total)). // limits to the exact number of rounds
		Group("location_id"). // ensures unique locations by grouping same location_id
		Column("location_id"). // select only the location id column
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return locations, nil
}