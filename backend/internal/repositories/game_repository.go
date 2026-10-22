package repositories

import (
	"backend/internal/models/games"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

var ErrNoLocationLeft = errors.New("no location left to replace a broken one")

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

// Marks the game as finished (finished_at), only once: calling it again changes nothing, so the hub can retry it
// safely. It does nothing else: the xp, rank and statistics of the players are updated separately (see the hub).
func (r *GameRepository) SetFinished(ctx context.Context,
	gameId int64) error {
	_, err := r.db.NewUpdate().
		Model((*games.Game)(nil)).
		Set("finished_at = ?", time.Now()).
		Where("id = ?", gameId).
		Where("finished_at IS NULL").
		Exec(ctx)
	return err
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

// Transaction that eplaces the location of the round roundNb of gameId by a random one
// that is neither used in this game nor disabled, and restarts the round with it:
//   - the old location is disabled so that it is never drawn again. It isn't deleted: the rounds already
//     played with it still point to it
//   - the new one takes its place in game.location_ids and in rounds.loc_id
//   - the round starts again from now (rounds.started_at), and its guesses (made without seeing anything) are deleted
// Returns the id of the old location and the id of the new one.
func (r *GameRepository) ReplaceRoundLocation(ctx context.Context, gameId, roundNb int64) (int64, int64, error) {
	var roundId, oldLocationId, newLocationId int64
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := tx.NewSelect().
			TableExpr("rounds").
			Column("id", "loc_id").
			Where("game_id = ? AND round_nb = ?", gameId, roundNb).
			Scan(ctx, &roundId, &oldLocationId); err != nil {
			return err
		}
		game := new(games.Game)
		if err := tx.NewSelect().
			Model(game).
			Column("location_ids").
			Where("id = ?", gameId).
			Scan(ctx); err != nil {
			return err
		}
		index := roundNb - 1 // because the array is 0-indexed, but round_nb starts at 1
		if index < 0 || index >= int64(len(game.LocationIDs)) {
			return fmt.Errorf("game %d has no location for round %d", gameId, roundNb)
		}
		// a random location, neither used in this game (the broken one included) nor disabled
		if err := tx.NewSelect().
			TableExpr("locations").
			Column("location_id").
			Where("disabled_at IS NULL").
			Where("location_id NOT IN (?)", bun.In(game.LocationIDs)).
			OrderExpr("RANDOM()").
			Limit(1).
			Scan(ctx, &newLocationId); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNoLocationLeft
			}
			return err
		}
		if _, err := tx.NewUpdate().
			TableExpr("locations").
			Set("disabled_at = ?", time.Now()).
			Where("location_id = ?", oldLocationId).
			Where("disabled_at IS NULL"). // already disabled by another game: it keeps its first date
			Exec(ctx); err != nil {
			return err
		}
		game.LocationIDs[index] = newLocationId
		if _, err := tx.NewUpdate().
			Model(game).
			Column("location_ids").
			Where("id = ?", gameId).
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewUpdate().
			TableExpr("rounds").
			Set("loc_id = ?", newLocationId).
			Set("started_at = ?", time.Now()).
			Where("id = ?", roundId).
			Exec(ctx); err != nil {
			return err
		}
		_, err := tx.NewDelete().
			TableExpr("guesses").
			Where("round_id = ?", roundId).
			Exec(ctx)
		return err
	})
	if err != nil {
		return 0, 0, err
	}
	return oldLocationId, newLocationId, nil
}

// Check if round_total is smaller than the number of available locations, the disabled ones excluded (else error)
func (r *LocationRepository) CheckEnoughLocations(ctx context.Context, round_total int64) (error) {
	var count int
	count, err := r.db.NewSelect().
		TableExpr("locations").
		Where("disabled_at IS NULL").
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
		Where("disabled_at IS NULL").
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