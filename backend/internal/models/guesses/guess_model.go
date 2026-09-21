package guesses

import(
	"time"
	"backend/internal/models/rounds"
	"backend/internal/models/users"
	"github.com/uptrace/bun"
)
// no injection allowed:
// json:"-" makes json.Unmarshal ignore them even if a client includes a field like "CountryFound", "Streak", etc. in the WS payload
type Guess struct {
	bun.BaseModel `bun:"table:guesses,alias:g"`

	ID				int64			`bun:"guess_id,pk,autoincrement"`
	Lat				*float64		`bun:"latitude" json:"latitude,omitempty"` // pointer to allow knowing if values have been set or not, so if user has guessed or not
	Long			*float64		`bun:"longitude" json:"longitude,omitempty"`
	CC				string			`bun:"country_code,type:varchar(2)" json:"-"`
	Distance		int64			`bun:"distance" json:"-"`
	CountryFound 	bool			`bun:"country_found,default:false" json:"-"`
	TimeLeft		time.Duration	`bun:"time_left,default:0" json:"-"` // duration are stored in nanoseconds, but can be shown in any type
	Streak			int64			`bun:"streak,default:0" json:"-"`
	BaseScore		int64			`bun:"base_score,notnull" json:"-"` // score avant multiplicateurs (streak, cc), utilisé pour l'XP
	Score			int64			`bun:"score,notnull" json:"-"` // score final avec multiplicateurs, utilisé pour le classement de la partie
	RoundNb 		int64			`bun:"round_nb,notnull" json:"round_nb,omitempty"`
	RoundID			int64			`bun:"round_id,notnull,unique:idx_guess_user_round" json:"-"`
	UserID			int64			`bun:"user_id,notnull,unique:idx_guess_user_round" json:"-"`

	Round			*rounds.Round	`bun:"rel:belongs-to,join:round_id=id"`
	User			*users.User		`bun:"rel:belongs-to,join:user_id=id"`
}