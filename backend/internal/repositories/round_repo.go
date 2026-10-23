package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"backend/internal/models/games"
	"backend/internal/models/rounds"
	"github.com/uptrace/bun"
)

var ErrRoundNotFound = errors.New("round not found")
var ErrGameNotFound = errors.New("game not found")

type RoundRepository struct {
	db                 *bun.DB
	locationRepository *LocationRepository
}

func NewRoundRepository(db *bun.DB, locationRepository *LocationRepository) *RoundRepository {
	return &RoundRepository{
		db:                 db,
		locationRepository: locationRepository,
	}
}

// CREATE - using round struct
func (r *RoundRepository) Create(ctx context.Context, round *rounds.Round) (*rounds.Round, error) {
	_, err := r.db.NewInsert().Model(round).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create round: %w", err)
	}
	return round, nil
}

// READ all - using round ID
func (r *RoundRepository) GetByID(ctx context.Context, round_id int64) (*rounds.Round, error) {
	round := new(rounds.Round)
	err := r.db.NewSelect().Model(round).Where("id = ?", round_id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoundNotFound
		}
		return nil, fmt.Errorf("failed to get round by ID: %w", err)
	}
	return round, nil
}

// READ - using game ID
func (r *RoundRepository) GetGame(ctx context.Context, game_id int64) error {
	var exists bool
	exists, err := r.db.NewSelect().Model((*games.Game)(nil)).Where("id = ?", game_id).Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if game exists: %w", err)
	}
	if !exists {
		return ErrGameNotFound
	}
	return nil
}

// READ - get the number of rounds already created for this game_id
func (r *RoundRepository) GetNextRoundNb(ctx context.Context, game_id int64) (int64, error) {
	var nextRound int64
	err := r.db.NewSelect().
		ColumnExpr("COALESCE(MAX(round_nb), 0) + 1 AS next_round").
		Table("rounds").
		Where("game_id = ?", game_id).
		Scan(ctx, &nextRound)
	if err != nil {
		return 0, fmt.Errorf("failed to count rounds for game_id %d: %w", game_id, err)
	}
	return nextRound, nil // return the next round number
}

// READ - check if the last round has been reached for this game_id
func (r *RoundRepository) CheckLastRound(ctx context.Context, game_id int64, round_nb int64) (bool, error) {
	var roundTotal int64
	err := r.db.NewSelect().
		Model((*games.Game)(nil)).
		Column("round_total").
		Where("id = ?", game_id).
		Limit(1).
		Scan(ctx, &roundTotal)
	if err != nil {
		return false, fmt.Errorf("failed to get round_total for game_id %d: %w", game_id, err)
	}
	if round_nb > roundTotal {
		return true, fmt.Errorf("failed for game_id %d: actual round %d exceeds total rounds %d", game_id, round_nb, roundTotal)
	}

	isLastRound := round_nb == roundTotal
	return isLastRound, nil
}

// READ - get location_id by using the map of locations associated with the game_id, and the round_nb (called id inside) as the key within that map
func (r *RoundRepository) GetLocationID(ctx context.Context, game_id int64, round_nb int64) (int64, error) {
	var location_id int64
	index := round_nb - 1 // because the array is 0-indexed, but round_nb starts at 1
	err := r.db.NewSelect().
		Model((*games.Game)(nil)).
		ColumnExpr("location_ids[?] AS location_id", index).
		Where("id = ?", game_id).
		Limit(1).
		Scan(ctx, &location_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrRoundNotFound
		}
		return 0, fmt.Errorf("failed to get location_id for game_id %d and round_nb %d: %w", game_id, round_nb, err)
	}
	return location_id, nil
}

func (r *RoundRepository) GetCoordinatesByLocationID(ctx context.Context, location_id int64) (float64, float64, error) {
	return r.locationRepository.GetCoordinatesByID(ctx, location_id)
}

func (r *RoundRepository) GetCCByLocationID(ctx context.Context, location_id int64) (string, error) {
	return r.locationRepository.GetCountryCodeByID(ctx, location_id)
}

func (r *RoundRepository) GetIDByGameIDAndRoundNb(ctx context.Context, game_id int64, round_nb int64) (int64, error) {
	var roundID int64
	err := r.db.NewSelect().
		Model((*rounds.Round)(nil)).
		Column("id").
		Where("game_id = ? AND round_nb = ?", game_id, round_nb).
		Scan(ctx, &roundID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrRoundNotFound
		}
		return 0, fmt.Errorf("failed to get round id by game id %d and round_nb %d: %w", game_id, round_nb, err)
	}
	return roundID, nil
}

// // DELETE - using round ID
// func (r *RoundRepository) DeleteByID(ctx context.Context, round_id int64) error {
// 	res, err := r.db.NewDelete().Model((*rounds.Round)(nil)).Where("round_id = ?", round_id).Exec(ctx)
// 	// Check for errors during the delete operation
// 	if err != nil {
// 		return fmt.Errorf("failed to delete round by ID: %w", err)
// 	}

// 	// Check if any rows were affected (i.e., if the round existed)
// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get rows affected after deleting round: %w", err)
// 	}
// 	// returns an error if no rows were affected, indicating that the round was not found
// 	if rowsAffected == 0 {
// 		return ErrRoundNotFound
// 	}

// 	return nil
// }
