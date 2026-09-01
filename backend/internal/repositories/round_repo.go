package repositories

import (
	"context"
	// "database/sql"
	"errors"
	"fmt"

	"backend/internal/models/rounds"
	"backend/internal/models/games"
	"github.com/uptrace/bun"
)

var ErrRoundNotFound = errors.New("round not found")
var ErrGameNotFound = errors.New("game not found")

type RoundRepository struct {
	db *bun.DB
}

func NewRoundRepository(db *bun.DB) *RoundRepository {
	return &RoundRepository{db: db}
}

// CREATE - using round struct
func (r *RoundRepository) Create(ctx context.Context, round *rounds.Round) (*rounds.Round, error) {
	_, err := r.db.NewInsert().Model(round).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create round: %w", err)
	}
	return round, nil
}

// // READ all - using round ID
// func (r *RoundRepository) GetByID(ctx context.Context, round_id int64) (*rounds.Round, error) {
// 	round := new(rounds.Round)
// 	err := r.db.NewSelect().Model(round).Where("round_id = ?", round_id).Scan(ctx)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, ErrRoundNotFound
// 		}
// 		return nil, fmt.Errorf("failed to get round by ID: %w", err)
// 	}
// 	return round, nil
// }

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

//READ - get the number of rounds already created for this game_id
func (r *RoundRepository) GetRoundNb(ctx context.Context, game_id int64) (int64, error) {
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

// // READ - get location_id by using the map of locations associated with the game_id, and the round_nb (called id inside) as the key within that map
// func (r *RoundRepository) GetLocationID(ctx context.Context, game_id int64, round_nb int64) (int64, error) {
// 	var location_id int64
// 	err := r.db.NewSelect().
//         Column("location_id").
//         Table("rounds").
//         Where("game_id = ? AND location_id = ?", game_id, round_nb). //location ou round_nb? check the schema > to get the right key
//         Limit(1).
//         Scan(ctx, &locationID)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return 0, ErrRoundNotFound
// 		}
// 		return 0, fmt.Errorf("failed to get location_id for game_id %d and round_nb %d: %w", game_id, round_nb, err)
// 	}
// 	return location_id, nil
// }

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