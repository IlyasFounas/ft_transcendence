package friends

import (
	"github.com/uptrace/bun"
)

type Friend struct {
	bun.BaseModel `bun:"table:friends,alias:f"`

	ID			int64	`bun:"id,pk,autoincrement"`
	User1_ID	int64	`bun:"user1_id,notnull"`
	User2_ID	int64	`bun:"user2_id,notnull"`
}