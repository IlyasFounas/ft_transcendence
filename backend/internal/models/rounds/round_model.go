package rounds

import(
	"github.com/uptrace/bun"
	// "backend/internal/models/locations"
	// "backend/internal/models/games"
)

type Round struct {
	bun.BaseModel `bun:"table:rounds,alias:r"`

	ID			int64					`bun:"id,pk,autoincrement"`
	RoundId		int64					`bun:"round_id,notnull"`
	// GameId		games.Game				`bun:"game_id,notnull,unique:idx_loc_game"` // TODO join + unique
	// LocationId	locations.Location.ID	`bun:"location_id,notnull,unique:idx_loc_game"` // TODO join + unique
	// StartedAt	time.Time				`bun:"started_at,notnull"`
	// Game			*games.Game				`bun:"rel:belongs-to,join:game_id=id"`
	// Location		*locations.Location		`bun:"rel:belongs-to,join:location_id=id"`
}