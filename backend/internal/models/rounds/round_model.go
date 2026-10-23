package rounds

import (
	"backend/internal/models/games"
	"backend/internal/models/locations"
	"github.com/uptrace/bun"
	"time"
)

type Round struct {
	bun.BaseModel `bun:"table:rounds,alias:r"`

	ID         int64               `bun:"id,pk,autoincrement" json:"id"`
	RoundNb    int64               `bun:"round_nb,notnull" json:"round_nb"`
	GameId     int64               `bun:"game_id,notnull" json:"game_id"`
	Game       *games.Game         `bun:"rel:belongs-to,join:game_id=id"`
	LocationId int64               `bun:"loc_id,notnull"`
	Location   *locations.Location `bun:"rel:belongs-to,join:loc_id=location_id"`
	StartedAt  time.Time           `bun:"started_at,notnull"`
}
