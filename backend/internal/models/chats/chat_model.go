package chats

import (
	"github.com/uptrace/bun"
	"time"
)

type Chat struct {
	bun.BaseModel `bun:"table:chats,alias:c"`
	
	ID				int64		`bun:"id,pk,autoincrement"`
	Message			string		`bun:"message,notnull"`
	SRC_user_ID		int64		`bun:"src_user_id,notnull"`
	DST_user_ID		int64		`bun:"dst_user_id,notnull"`
	CreatedAt		time.Time	`bun:"created_at,default:current_timestamp"`
}