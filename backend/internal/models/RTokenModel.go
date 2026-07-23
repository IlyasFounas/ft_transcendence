package model

import (
    "github.com/uptrace/bun"
	"time"
)

type RefreshToken struct {
    bun.BaseModel `bun:"table:refresh_tokens,alias:at"`

    ID			int64      `bun:"id,pk,autoincrement"`
    HashToken	string     `bun:"hash_token,notnull,unique"`
    UserID		int64       `bun:"user_id,notnull"`
    ExpiresAt	time.Time  `bun:"expires_at,notnull"`

    User *User `bun:"rel:belongs-to,join:user_id=id"`
}