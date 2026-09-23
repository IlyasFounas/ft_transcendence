package rounds

import(
	"time"
	"github.com/uptrace/bun"
	"backend/internal/models/locations"
	"backend/internal/models/games"
)

type Round struct {
	bun.BaseModel `bun:"table:rounds,alias:r"`

	ID			int64					`bun:"id,pk,autoincrement" json:"id"`
	RoundNb		int64					`bun:"round_nb,notnull" json:"round_nb"`
	LastRound	bool					`bun:"last_round,notnull,default:false"`
	GameId		int64					`bun:"game_id,notnull" json:"game_id"`
	Game		*games.Game				`bun:"rel:belongs-to,join:game_id=id"`
	LocationId	int64					`bun:"loc_id,notnull"` // TODO join + unique ,unique:idx_loc_game
	Location	*locations.Location		`bun:"rel:belongs-to,join:loc_id=location_id"`
	// StartedAt is the instant this round's GUESS PHASE begins — i.e. once its fixed pre-round
	// pause (see ws_handlers.prePlayPause) has elapsed, not the instant the row was written. The
	// row is created as soon as the pause starts, with StartedAt already set to that future
	// instant, so StartedAt alone (compared to time.Now(), and to Game.Timer for the round's
	// duration) is enough to tell whether the round is still pending, currently playing, or
	// already over — see RoundWsHandler.BeginRound / recoverRoundState.
	//
	// This round's duration is intentionally NOT a field here: it always comes from Game.Timer,
	// read live, since Game is the single source of truth for how long every round of a game
	// lasts.
	StartedAt	time.Time				`bun:"started_at,notnull"`
}