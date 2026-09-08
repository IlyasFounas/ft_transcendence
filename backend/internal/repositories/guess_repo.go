package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"backend/internal/models/guesses"
	"backend/internal/models/rounds"
	// "backend/internal/models/locations"
	"github.com/uptrace/bun"
)

var ErrGuessNotFound = errors.New("guess not found")

type GuessRepository struct {
	db *bun.DB
}

func NewGuessRepository(db *bun.DB) *GuessRepository {
	return &GuessRepository{db: db}
}

// CREATE - using guess struct
func (r *GuessRepository) Create(ctx context.Context, guess *guesses.Guess) (*guesses.Guess, error) {
	_, err := r.db.NewInsert().Model(guess).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create guess: %w", err)
	}
	return guess, nil
}

// READ all - using guess ID
func (r *GuessRepository) GetByID(ctx context.Context, guess_id int64) (*guesses.Guess, error) {
	guess := new(guesses.Guess)
	err := r.db.NewSelect().Model(guess).Where("guess_id = ?", guess_id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGuessNotFound
		}
		return nil, fmt.Errorf("failed to get guess by ID: %w", err)
	}
	return guess, nil
}

// READ coordinates - using guess ID
func (r *GuessRepository) GetCoordinatesByID(ctx context.Context, guess_id int64) (*guesses.Guess, error) {
	coordinates := new(guesses.Guess)
	err := r.db.NewSelect().Model(coordinates).Column("latitude", "longitude").Where("guess_id = ?", guess_id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGuessNotFound
		}
		return nil, fmt.Errorf("failed to get guess coordinates by ID: %w", err)
	}
	return coordinates, nil
}

func (r *GuessRepository) GetLocationIDByRoundID(ctx context.Context, round_id int64) (int64, error) {
	var locID int64
	err := r.db.NewSelect().Model((*rounds.Round)(nil)).Column("loc_id").Where("id = ?", round_id).Limit(1).Scan(ctx, &locID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrGuessNotFound
		}
		return 0, fmt.Errorf("failed to get location id by round id: %w", err)
	}
	return locID, nil
}

// func (r *GuessRepository) GetCountryCodeByLocationID(ctx context.Context, location_id int64) (string, error) {
// 	var cc string
// 	err := r.db.NewSelect().Model((*locations.Location)(nil)).Column("country_code").Where("location_id = ?", location_id).Limit(1).Scan(ctx, &cc)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return "", ErrGuessNotFound
// 		}
// 		return "", fmt.Errorf("failed to get country code by location id: %w", err)
// 	}
// 	return cc, nil
// }

// READ latitude - using guess ID
func (r *GuessRepository) GetLatitudeByID(ctx context.Context, guess_id int64) (float64, error) {
	var latitude float64
	err := r.db.NewSelect().Model((*guesses.Guess)(nil)).Column("latitude").Where("guess_id = ?", guess_id).Scan(ctx, &latitude)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrGuessNotFound
		}
		return 0, fmt.Errorf("failed to get guess latitude by ID: %w", err)
	}
	return latitude, nil
}

// READ longitude - using guess ID
func (r *GuessRepository) GetLongitudeByID(ctx context.Context, guess_id int64) (float64, error) {
	var longitude float64
	err := r.db.NewSelect().Model((*guesses.Guess)(nil)).Column("longitude").Where("guess_id = ?", guess_id).Scan(ctx, &longitude)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrGuessNotFound
		}
		return 0, fmt.Errorf("failed to get guess longitude by ID: %w", err)
	}
	return longitude, nil
}

// READ country code - using guess ID
func (r *GuessRepository) GetCountryCodeByID(ctx context.Context, guess_id int64) (string, error) {
	var countryCode string
	err := r.db.NewSelect().Model((*guesses.Guess)(nil)).Column("country_code").Where("guess_id = ?", guess_id).Scan(ctx, &countryCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrGuessNotFound
		}
		return "", fmt.Errorf("failed to get guess country code by ID: %w", err)
	}
	return countryCode, nil
}

// READ distance - using guess ID
func (r *GuessRepository) GetDistanceByID(ctx context.Context, guess_id int64) (int64, error) {
	var distance int64
	err := r.db.NewSelect().Model((*guesses.Guess)(nil)).Column("distance").Where("guess_id = ?", guess_id).Scan(ctx, &distance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrGuessNotFound
		}
		return 0, fmt.Errorf("failed to get guess distance by ID: %w", err)
	}
	return distance, nil
}

// READ score - using guess ID
func (r *GuessRepository) GetScoreByID(ctx context.Context, guess_id int64) (int64, error) {
	var score int64
	err := r.db.NewSelect().Model((*guesses.Guess)(nil)).Column("score").Where("guess_id = ?", guess_id).Scan(ctx, &score)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrGuessNotFound
		}
		return 0, fmt.Errorf("failed to get guess score by ID: %w", err)
	}
	return score, nil
}

// DELETE - using guess ID
func (r *GuessRepository) DeleteByID(ctx context.Context, guess_id int64) error {
	res, err := r.db.NewDelete().Model((*guesses.Guess)(nil)).Where("guess_id = ?", guess_id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete guess by ID: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected after deleting guess: %w", err)
	}

	if rowsAffected == 0 {
		return ErrGuessNotFound
	}

	return nil
}
