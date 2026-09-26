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
	RoundTotal 	int64 		`bun:"round_total,notnull"`
	FinishedAt	*time.Time 	`bun:"finished_at"`
	PlayerNb	int64 		`bun:"player_nb,notnull,type:bigint CHECK (player_nb <> 0)"`
	CreatedAt   time.Time	`bun:"created_at,default:current_timestamp"`

	LocationIDs []int64		`bun:"location_ids,type:jsonb" json:"location_ids"`
	Lobby 		[]*Lobby  	`bun:"rel:has-many,join:id=game_id"`
}

type Lobby struct {
	bun.BaseModel `bun:"table:lobby,alias:lo"`

	ID 			int64		`bun:"id,pk,autoincrement"`
	UserId  	int64		`bun:"user_id"`
	GameId  	int64		`bun:"game_id"`
	JoinCode  	*string		`bun:"join_code,unique"`
	IsCreator  	bool 		`bun:"is_creator"`
	LeftAt  	*time.Time 	`bun:"left_at"`
	Kicked  	bool 		`bun:"kicked,notnull,default:false" json:"-"` // kicked out by the host: can't join this game again

	Game 		*Game 	`bun:"rel:has-one,join:game_id=id"`
}