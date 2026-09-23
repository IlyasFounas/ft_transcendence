package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"backend/internal/models/guesses"
	"backend/internal/models/rounds"

	"github.com/uptrace/bun"
)

var ErrGuessNotFound = errors.New("guess not found")

type LeaderboardEntry struct {
	Position int64  `json:"position"`
	Username string `json:"username"`
	Total    int64  `json:"total"`
}

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

// "backend/internal/models/locations"
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

// READ - using round ID (to get the score for that user in that round)
func (r *GuessRepository) GetScoreByRoundAndUserID(ctx context.Context, round_id int64, user_id int64) (int64, error) {
	var score int64
	err := r.db.NewSelect().
		Model((*guesses.Guess)(nil)).
		Column("score").
		Where("round_id = ? AND user_id = ?", round_id, user_id).
		Scan(ctx, &score)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrGuessNotFound
		}
		return 0, fmt.Errorf("failed to get score for round %d and user %d: %w", round_id, user_id, err)
	}
	return score, nil
}

// READ - using game ID and user ID (to get the total score for that user in that game)
func (r *GuessRepository) GetTotalScoreByGameAndUserID(ctx context.Context, game_id int64, user_id int64) (int64, error) {
	var total int64
	err := r.db.NewSelect().
		Model((*guesses.Guess)(nil)).
		ColumnExpr("COALESCE(SUM(g.score), 0)"). // COALESCE(..., 0) évite un NULL si l'utilisateur n'a encore aucun guess pour cette game (au lieu de creer une erreur)
		Join("JOIN rounds AS r ON r.id = g.round_id").
		Where("r.game_id = ? AND g.user_id = ?", game_id, user_id).
		Scan(ctx, &total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total score for game %d and user %d: %w", game_id, user_id, err)
	}
	return total, nil
}

// READ - using user ID (to get the XP for that user across all games)
func (r *GuessRepository) GetXPByUserID(ctx context.Context, user_id int64) (int64, error) {
	var total int64
	err := r.db.NewSelect().
		Model((*guesses.Guess)(nil)).
		ColumnExpr("COALESCE(SUM(base_score), 0)").
		Where("user_id = ?", user_id).
		Scan(ctx, &total)
	if err != nil {
		return 0, fmt.Errorf("failed to get XP for user %d: %w", user_id, err)
	}
	return total, nil
}

// READ - using game ID (to get the results for that game)
func (r *GuessRepository) GetResultsByGameID(ctx context.Context, game_id int64) ([][2]int64, error) {
	type results struct {
		UserID int64 `bun:"user_id"`
		Total  int64 `bun:"total"`
	}

	var rows []results
	err := r.db.NewSelect().
		Model((*guesses.Guess)(nil)).
		ColumnExpr("g.user_id AS user_id").
		ColumnExpr("SUM(g.score) AS total").
		Join("JOIN rounds AS r ON r.id = g.round_id").
		Where("r.game_id = ?", game_id).
		Group("g.user_id").
		OrderExpr("total DESC").
		Scan(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("failed to get results for game %d: %w", game_id, err)
	}
	res := make([][2]int64, 0, len(rows))
	for _, row := range rows {
		res = append(res, [2]int64{row.UserID, row.Total})
	}

	return res, nil
}

// READ - using game ID (to get the leaderboard for all games: 100 first players based on position, usernames and total xp, and position, username and total xp of user who requested the leaderboard)
func (r *GuessRepository) GetLeaderboard(ctx context.Context, userID int64) ([]LeaderboardEntry, error) {
	type leaderboard struct {
		Username string `bun:"username"`
		Total    int64  `bun:"total"`
		Position int64  `bun:"position"`
	}

	// Get all users ranked by total score (including users with 0 score)
	var rows []leaderboard
	err := r.db.NewSelect().
		TableExpr("users AS u").
		ColumnExpr("u.name AS username").
		ColumnExpr("COALESCE(SUM(g.base_score), 0) AS total").
		ColumnExpr("RANK() OVER (ORDER BY COALESCE(SUM(g.base_score), 0) DESC) AS position").
		Join("LEFT JOIN guesses AS g ON g.user_id = u.id").
		Group("u.id").
		Group("u.name").
		OrderExpr("total DESC").
		Limit(100).
		Scan(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	// Retrieve position, username and total xp for the requesting user (ensures they appear even with 0 score)
	var userRow leaderboard
	errUser := r.db.NewSelect().
		TableExpr("users AS u").
		ColumnExpr("u.name AS username").
		ColumnExpr("COALESCE(SUM(g.base_score), 0) AS total").
		ColumnExpr("RANK() OVER (ORDER BY COALESCE(SUM(g.base_score), 0) DESC) AS position").
		Join("LEFT JOIN guesses AS g ON g.user_id = u.id").
		Where("u.id = ?", userID).
		Group("u.id").
		Group("u.name").
		Scan(ctx, &userRow)

	// Check if user is already in rows (to avoid duplicates)
	isUserInRows := false
	for _, row := range rows {
		if row.Username == userRow.Username {
			isUserInRows = true
			break
		}
	}

	// Add user row if not already in top 100 and query was successful
	if errUser == nil && !isUserInRows {
		rows = append(rows, userRow)
	}

	res := make([]LeaderboardEntry, 0, len(rows))
	for _, row := range rows {
		res = append(res, LeaderboardEntry{
			Position: row.Position,
			Username: row.Username,
			Total:    row.Total,
		})
	}

	return res, nil
}

// READ - streak value for the last round - using game ID and user ID
func (r *GuessRepository) GetLastStreakByGameAndUserID(ctx context.Context, game_id int64, user_id int64, round_nb int64) (int64, error) {
	if round_nb <= 1 {
		return 0, nil
	}
	var streak int64
	err := r.db.NewSelect().
		Model((*guesses.Guess)(nil)).
		ColumnExpr("g.streak").
		Join("JOIN rounds AS r ON r.id = g.round_id").
		Where("r.game_id = ? AND g.user_id = ? AND r.round_nb = ?", game_id, user_id, round_nb-1).
		Limit(1).
		Scan(ctx, &streak)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get last streak for user %d in game %d: %w", user_id, game_id, err)
	}
	return streak, nil
}

// READ - using round ID (to get all guesses for that round)
func (r *GuessRepository) GetByRoundID(ctx context.Context, round_id int64) ([]*guesses.Guess, error) {
	var guesses []*guesses.Guess
	err := r.db.NewSelect().
		Model(&guesses).
		Where("round_id = ?", round_id).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get guesses by round ID: %w", err)
	}
	return guesses, nil
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
