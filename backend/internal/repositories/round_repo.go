package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"backend/internal/models/rounds"
	"github.com/uptrace/bun"
)

var ErrRoundNotFound = errors.New("round not found")

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

// READ all - using round ID
func (r *RoundRepository) GetByID(ctx context.Context, round_id int64) (*rounds.Round, error) {
	round := new(rounds.Round)
	err := r.db.NewSelect().Model(round).Where("round_id = ?", round_id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoundNotFound
		}
		return nil, fmt.Errorf("failed to get round by ID: %w", err)
	}
	return round, nil
}

// DELETE - using round ID
func (r *RoundRepository) DeleteByID(ctx context.Context, round_id int64) error {
	res, err := r.db.NewDelete().Model((*rounds.Round)(nil)).Where("round_id = ?", round_id).Exec(ctx)
	// Check for errors during the delete operation
	if err != nil {
		return fmt.Errorf("failed to delete round by ID: %w", err)
	}

	// Check if any rows were affected (i.e., if the round existed)
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected after deleting round: %w", err)
	}
	// returns an error if no rows were affected, indicating that the round was not found
	if rowsAffected == 0 {
		return ErrRoundNotFound
	}

	return nil
}