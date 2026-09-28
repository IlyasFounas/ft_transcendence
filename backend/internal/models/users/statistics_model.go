package users

import (
	"github.com/uptrace/bun"
)

type Statistics struct {
	bun.BaseModel `bun:"table:statistics,alias:st"`

	ID            int64   `bun:"id,pk,autoincrement"`
	UserID        int64   `bun:"user_id,notnull,unique"`
	FinishedGames int64   `bun:"finished_games,notnull,default:0"` // games he was still in when they ended
	Accuracy      float64 `bun:"accuracy,notnull,default:0"`       // percentage (0-100): average accuracy of his finished games
	BestScore     int64   `bun:"best_score,notnull,default:0"`     // best final score of a game (score with multipliers, at the very end)
	ClosestGuess  *int64  `bun:"closest_guess"`                    // closest distance he ever guessed, in meters (nil: no guess yet)

	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

type Distance struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// Converts a distance in meters into a Distance ready to display
func NewDistance(meters int64) Distance {
	if meters < 1000 {
		return Distance{Value: float64(meters), Unit: "m"}
	}
	return Distance{Value: float64(meters) / 1000, Unit: "km"}
}