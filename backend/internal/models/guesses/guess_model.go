package guesses

import(
	// "time"
	"github.com/uptrace/bun"
	// "backend/internal/models/games"
	// "backend/internal/models/rounds"
)

type Guess struct {
	bun.BaseModel `bun:"table:guesses,alias:g"`

	ID				int64			`bun:"guess_id,pk,autoincrement"`
	Lat				*float64		`bun:"latitude" json:"latitude,omitempty"` // pointer to allow knowing if values have been set or not, so if user has guessed or not
	Long			*float64		`bun:"longitude" json:"longitude,omitempty"`
	CC				string			`bun:"country_code,type:varchar(2)"`
	Distance		int64			`bun:"distance"`
	CountryFound 	bool			`bun:"country_found,default:false"`
	// TimeTaken		time.Time		`bun:"time_taken,default:0"`
	Streak			int64			`bun:"streak,default:0"`
	Score			int64			`bun:"score,notnull"`
	// RoundId			int64			`bun:"round_id,notnull,unique"`
	// UserId			int64			`bun:"user_id,notnull, unique"`

	// Round			*rounds.Round	`bun:"rel:belongs-to,join:round_id=id"`
}