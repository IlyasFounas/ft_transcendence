package users

import (
	"github.com/uptrace/bun"
)

type Statistics struct {
	bun.BaseModel `bun:"table:statistics,alias:st"`

	ID				int64	`bun:"id,pk,autoincrement"`
	UserID			int64	`bun:"user_id,notnull,unique"`
	FinishedGames	int64	`bun:"finished_games,notnull,default:0" json:"finished_games, notnull"`
	Accuracy		float64	`bun:"accuracy,notnull,default:0"`
	BestScore		int64	`bun:"best_score,notnull,default:0"`
	ClosestGuess	*int64	`bun:"closest_guess"`

	User			*User	`bun:"rel:belongs-to,join:user_id=id"`
}

type Distance struct {
	Value	float64	`json:"value"`
	Unit	string	`json:"unit"`
}

// Converts a distance in meters into a Distance ready to display
func NewDistance(meters int64) Distance {
	if meters < 1000 {
		return Distance{Value: float64(meters), Unit: "m"}
	}
	return Distance{Value: float64(meters) / 1000, Unit: "km"}
}