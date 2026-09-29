package repositories

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/models/users"
	"backend/internal/tools/haversine"

	"github.com/uptrace/bun"
)

type StatisticsRepository struct {
	db *bun.DB
}

func NewStatisticsRepository(db *bun.DB) *StatisticsRepository {
	return &StatisticsRepository{db: db}
}

// Statistics returns every statistic of userId
// A user who has never finished a game has no row yet: he gets zero statistics, and no closest guess.
func (r *StatisticsRepository) Statistics(ctx context.Context, userId int64) (*users.Statistics, error) {
	stats := new(users.Statistics)
	err := r.db.NewSelect().
		Model(stats).
		Where("user_id = ?", userId).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return &users.Statistics{UserID: userId}, nil
	}
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// Getters
func (r *StatisticsRepository) FinishedGames(ctx context.Context, userId int64) (int64, error) {
	stats, err := r.Statistics(ctx, userId)
	if err != nil {
		return 0, err
	}
	return stats.FinishedGames, nil
}

func (r *StatisticsRepository) Accuracy(ctx context.Context, userId int64) (float64, error) {
	stats, err := r.Statistics(ctx, userId)
	if err != nil {
		return 0, err
	}
	return stats.Accuracy, nil
}

func (r *StatisticsRepository) BestScore(ctx context.Context, userId int64) (int64, error) {
	stats, err := r.Statistics(ctx, userId)
	if err != nil {
		return 0, err
	}
	return stats.BestScore, nil
}

// Returns the closest distance ever guessed by userId, in meters (nil: no guess yet)
func (r *StatisticsRepository) ClosestGuess(ctx context.Context, userId int64) (*int64, error) {
	stats, err := r.Statistics(ctx, userId)
	if err != nil {
		return nil, err
	}
	return stats.ClosestGuess, nil
}

// Setters
func (r *StatisticsRepository) SetFinishedGames(ctx context.Context, userId int64, finishedGames int64) (bool, error) {
	if finishedGames < 0 {
		return false, nil
	}
	if err := r.set(ctx, &users.Statistics{UserID: userId, FinishedGames: finishedGames}, "finished_games"); err != nil {
		return false, err
	}
	return true, nil
}

func (r *StatisticsRepository) SetAccuracy(ctx context.Context, userId int64, accuracy float64) (bool, error) {
	if accuracy < 0 || accuracy > 100 {
		return false, nil
	}
	if err := r.set(ctx, &users.Statistics{UserID: userId, Accuracy: accuracy}, "accuracy"); err != nil {
		return false, err
	}
	return true, nil
}

func (r *StatisticsRepository) SetBestScore(ctx context.Context, userId int64, bestScore int64) (bool, error) {
	if bestScore < 0 {
		return false, nil
	}
	if err := r.set(ctx, &users.Statistics{UserID: userId, BestScore: bestScore}, "best_score"); err != nil {
		return false, err
	}
	return true, nil
}

// Takes a distance in meters
func (r *StatisticsRepository) SetClosestGuess(ctx context.Context, userId int64, meters int64) (bool, error) {
	if meters < 0 {
		return false, nil
	}
	if err := r.set(ctx, &users.Statistics{UserID: userId, ClosestGuess: &meters}, "closest_guess"); err != nil {
		return false, err
	}
	return true, nil
}

// Writes one column of the statistics of stats.UserID: it inserts his row if he has none yet (the other values
// of stats are zero), otherwise it only replaces this column
func (r *StatisticsRepository) set(ctx context.Context, stats *users.Statistics, column string) error {
	_, err := r.db.NewInsert().
		Model(stats).
		On("CONFLICT (user_id) DO UPDATE").
		Set("? = EXCLUDED.?", bun.Ident(column), bun.Ident(column)).
		Exec(ctx)
	return err
}

// Updatethe statistics of each user from all his FINISHED games
//   - finished_games nb (a game without any guess of his counts too)
//   - accuracy: average over these games of the accuracy of the game, which is the average over its rounds of the
//     accuracy of his guess: a rule of three on the distance (0 m = 100 %, haversine.EarthDiameterMeters or more
//     = 0 %), a round without guess counting 0 %
//   - best_score: best final score of a game (sum of his scores with multipliers)
//   - closest_guess: closest distance the player ever guessed, in meters (or km if more than 1000 m), nil if he never guessed
func (r *StatisticsRepository) UpdateStats(ctx context.Context, userIds []int64) error {
	if len(userIds) == 0 {
		return nil
	}

	// number of finished games of each user
	var counts []struct {
		UserID int64 `bun:"user_id"`
		Games  int64 `bun:"games"`
	}
	err := finishedGamesLobbies(r.db).
		ColumnExpr("lo.user_id").
		ColumnExpr("COUNT(DISTINCT lo.game_id) AS games").
		Where("lo.user_id IN (?)", bun.In(userIds)).
		Group("lo.user_id").
		Scan(ctx, &counts)
	if err != nil {
		return err
	}

	// what each user did in each of his finished games
	var perGame []struct {
		UserID   int64   `bun:"user_id"`
		Score    int64   `bun:"score"`    // final score of the game
		Accuracy float64 `bun:"accuracy"` // accuracy of the game (average over its rounds)
		Closest  *int64  `bun:"closest"`  // closest guess of the game, in meters
	}
	err = finishedGamesGuesses(r.db).
		ColumnExpr("g.user_id").
		ColumnExpr("SUM(g.score) AS score").
		ColumnExpr("SUM(GREATEST(0, 100 - g.distance * 100.0 / ?)) / MAX(ga.round_total) AS accuracy",
			haversine.EarthDiameterMeters).
		ColumnExpr("MIN(g.distance) AS closest").
		Where("g.user_id IN (?)", bun.In(userIds)).
		GroupExpr("g.user_id, r.game_id").
		Scan(ctx, &perGame)
	if err != nil {
		return err
	}

	// statistics of each user, from scratch
	stats := make(map[int64]*users.Statistics, len(userIds))
	for _, userId := range userIds {
		stats[userId] = &users.Statistics{UserID: userId}
	}
	for _, count := range counts {
		stats[count.UserID].FinishedGames = count.Games
	}
	for _, game := range perGame {
		s := stats[game.UserID]
		s.Accuracy += game.Accuracy // summed here, divided by the number of finished games below
		if game.Score > s.BestScore {
			s.BestScore = game.Score
		}
		if game.Closest != nil && (s.ClosestGuess == nil || *game.Closest < *s.ClosestGuess) {
			closest := *game.Closest
			s.ClosestGuess = &closest
		}
	}

	// every value is replaced (nothing is merged with the previous row). Every user is tried even if one fails.
	var errs []error
	for _, s := range stats {
		if s.FinishedGames > 0 {
			s.Accuracy /= float64(s.FinishedGames) // games without any guess count 0 %
		}
		_, err := r.db.NewInsert().
			Model(s).
			On("CONFLICT (user_id) DO UPDATE").
			Set("finished_games = EXCLUDED.finished_games").
			Set("accuracy = EXCLUDED.accuracy").
			Set("best_score = EXCLUDED.best_score").
			Set("closest_guess = EXCLUDED.closest_guess").
			Exec(ctx)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}