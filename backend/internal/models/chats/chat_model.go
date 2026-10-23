package chats

import (
	"github.com/uptrace/bun"
	"time"
)

type Chat struct {
	bun.BaseModel `bun:"table:chats,alias:c"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	SRC_user_ID int64     `bun:"src_user_id,notnull" json:"user_id"`
	DST_user_ID int64     `bun:"dst_user_id,notnull" json:"friend_id"`
	Message     string    `bun:"message,notnull" json:"message"`
	CreatedAt   time.Time `bun:"created_at,default:current_timestamp" json:"date"`
}

type Cursor struct {
	Start *int64 `json:"start,omitempty"`
	End   *int64 `json:"end,omitempty"`
}

func NewCursor(chats []*Chat) Cursor {
	cursor := Cursor{}
	if len(chats) > 0 {
		cursor.Start = &chats[0].ID
		cursor.End = &chats[len(chats)-1].ID
	}
	return cursor
}
