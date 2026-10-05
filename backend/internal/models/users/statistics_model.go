package users

import (
	"math"
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

func roundToDecimals(value float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(value *ratio) / ratio
}

// Converts and rounds a distance in meters into a Distance ready to display
func NewDistance(meters int64) Distance {
	if meters < 10000 {
		return Distance{Value: float64(meters), Unit: "m"}
	}
	km := float64(meters) / 1000
	if km < 100 {
		return Distance{Value: roundToDecimals(km, 1), Unit: "km"}
	} else {
		return Distance{Value: roundToDecimals(km, 0), Unit: "km"}
	}
}