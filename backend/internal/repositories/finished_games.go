package repositories

import (
	"backend/internal/models/games"
	"backend/internal/models/guesses"

	"github.com/uptrace/bun"
)

// The games that count for a player (xp, rank, statistics) are his FINISHED games: the games that are over
// (finished_at set) and that he was still in when they ended, that is not left, or left after finished_at.
// A game he left or was kicked out of before its end doesn't count for him.

// stillInAtTheEnd is the condition "the player of the lobby row lo was still in the game ga when it ended"
func stillInAtTheEnd(q *bun.SelectQuery) *bun.SelectQuery {
	return q.Where("lo.left_at IS NULL").WhereOr("lo.left_at >= ga.finished_at")
}

// finishedGamesLobbies selects the lobby rows of the finished games of the players (alias lo, with ga = the game):
// one per finished game, even if the player didn't guess once in it. The caller adds its columns and conditions.
func finishedGamesLobbies(db bun.IDB) *bun.SelectQuery {
	return db.NewSelect().
		Model((*games.Lobby)(nil)).
		Join("JOIN game AS ga ON ga.id = lo.game_id").
		Where("ga.finished_at IS NOT NULL").
		WhereGroup(" AND ", stillInAtTheEnd)
}

// finishedGamesGuesses selects the guesses of the players in their finished games (alias g, with r = its round and
// ga = its game). The caller adds its columns, conditions and grouping.
func finishedGamesGuesses(db bun.IDB) *bun.SelectQuery {
	inTheGameAtTheEnd := db.NewSelect().
		Model((*games.Lobby)(nil)).
		ColumnExpr("1").
		Where("lo.game_id = ga.id").
		Where("lo.user_id = g.user_id").
		WhereGroup(" AND ", stillInAtTheEnd)
	return db.NewSelect().
		Model((*guesses.Guess)(nil)).
		Join("JOIN rounds AS r ON r.id = g.round_id").
		Join("JOIN game AS ga ON ga.id = r.game_id").
		Where("ga.finished_at IS NOT NULL").
		Where("EXISTS (?)", inTheGameAtTheEnd)
}