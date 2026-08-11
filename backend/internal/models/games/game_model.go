package games

import (
	"time"

	"github.com/uptrace/bun"
)

type Game struct {
	bun.BaseModel `bun:"table:game,alias:g"`

	ID			int64 		`bun:"id,pk,autoincrement"`
	Title		string 		`bun:"title"`
	Mode		string 		`bun:"mode"`
	Difficulty	string 		`bun:"difficulty"`
	Timer		int64 		`bun:"timer,notnull"`
	RoundNb 	int64 		`bun:"round_nb,notnull,check:round_nb <> 0"`
	Finished	bool 		`bun:"finished"`
	PlayerNb	int64 		`bun:"player_nb,notnull,check player_nb <> 0`
	CreatedAt   time.Time	`bun:"created_at,default:current_timestamp"`
}