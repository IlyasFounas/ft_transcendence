package friends

import (
	"github.com/uptrace/bun"
	"backend/internal/models/users"
	"time"
)

type Friend struct {
	bun.BaseModel `bun:"table:friends,alias:f"`

	ID				int64		`bun:"id,pk,autoincrement"`
	Requester_ID	int64		`bun:"requester_id,notnull"`
	Receiver_ID		int64		`bun:"receiver_id,notnull"`
	Status			string		`bun:"status,notnull,default:'pending'"`
	CreatedAt		time.Time	`bun:"created_at,default:current_timestamp"`
	
	Requester 		*users.User `bun:"rel:belongs-to,join:requester_id=id"`
}