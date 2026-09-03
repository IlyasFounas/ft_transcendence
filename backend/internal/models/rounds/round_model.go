package rounds

import(
	"github.com/uptrace/bun"
	// "backend/internal/models/locations"
	"backend/internal/models/games"
)

type Round struct {
	bun.BaseModel `bun:"table:rounds,alias:r"`

	ID			int64					`bun:"id,pk,autoincrement" json:"id"`
	RoundNb		int64					`bun:"round_nb,notnull" json:"round_nb"`
	LastRound	bool					`bun:"last_round,notnull,default:false"`
	GameId		int64					`bun:"game_id,notnull" json:"game_id"` // TODO join + unique ,unique:idx_loc_game
	Game		*games.Game				`bun:"rel:belongs-to,join:game_id=id"`
	// LocationId	int64					`bun:"location_id,notnull"` // TODO join + unique ,unique:idx_loc_game
	// Location	*locations.Location		`bun:"rel:belongs-to,join:location_id=id"`
	// StartedAt	time.Time				`bun:"started_at,notnull"`
}