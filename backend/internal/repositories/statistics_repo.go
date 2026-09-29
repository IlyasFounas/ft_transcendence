package repositories

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/models/games"
	"backend/internal/models/guesses"
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

// gameStatistics is what one player did in one game
type gameStatistics struct {
	UserID   int64   `bun:"user_id"`
	Score    int64   `bun:"score"`    // final score of the game: sum of the scores of his guesses (with multipliers)
	Accuracy float64 `bun:"accuracy"` // sum of the accuracy of his guesses
	Closest  *int64  `bun:"closest"`  // closest distance of his guesses, in meters
}

func (r *StatisticsRepository) GetRoundTotal(ctx context.Context, gameId int64) (int64, error) {
	var roundTotal int64
	err := r.db.NewSelect().
		Model((*games.Game)(nil)).
		Column("round_total").
		Where("id = ?", gameId).
		Scan(ctx, &roundTotal)
	if err != nil {
		return 0, err
	}
	return roundTotal, nil
}

func (r *StatisticsRepository) GetPlayersInGame(ctx context.Context, gameId int64) ([]*games.Lobby, error) {
	var players []*games.Lobby
	err := r.db.NewSelect().
		Model(&players).
		Where("game_id = ?", gameId).
		Where("left_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return players, nil
}

func (r *StatisticsRepository) SetGameStatisticsPerUser(ctx context.Context, gameId int64) (map[int64]gameStatistics, error) {
		var rows []gameStatistics
	err := r.db.NewSelect().
		Model((*guesses.Guess)(nil)).
		ColumnExpr("g.user_id").
		ColumnExpr("SUM(g.score) AS score"). // 
		ColumnExpr("SUM(GREATEST(0, 100 - g.distance * 100.0 / ?)) AS accuracy", haversine.EarthDiameterMeters).
		ColumnExpr("MIN(g.distance) AS closest").
		Join("JOIN rounds AS r ON r.id = g.round_id").
		Where("r.game_id = ?", gameId).
		Group("g.user_id").
		Scan(ctx, &rows)
	if err != nil {
		return nil, err
	}
	byUser := make(map[int64]gameStatistics, len(rows))
	for _, row := range rows {
		byUser[row.UserID] = row
	}
	return byUser, nil
}

func (r *StatisticsRepository) UpdateStatistics(ctx context.Context, userGameStats map[int64]gameStatistics, players []*games.Lobby, roundTotal int64) error {
	for _, player := range players {
		game := userGameStats[player.UserId] // zero value if he didn't guess once in this game
		stats := &users.Statistics{
			UserID:        player.UserId,
			FinishedGames: 1,
			Accuracy:      game.Accuracy / float64(roundTotal),
			BestScore:     game.Score,
			ClosestGuess:  game.Closest,
		}
		// first finished game: the row is inserted with these values. Otherwise they are merged with his row
		// (st = his current values, EXCLUDED = this game; LEAST ignores a NULL closest guess)
		_, err := r.db.NewInsert().
			Model(stats).
			On("CONFLICT (user_id) DO UPDATE").
			Set("finished_games = st.finished_games + 1").
			Set("accuracy = (st.accuracy * st.finished_games + EXCLUDED.accuracy) / (st.finished_games + 1)").
			Set("best_score = GREATEST(st.best_score, EXCLUDED.best_score)").
			Set("closest_guess = LEAST(st.closest_guess, EXCLUDED.closest_guess)").
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}