package tokens

import (
	"time"
    "github.com/uptrace/bun"
    "backend/internal/models/users"
)

type RefreshToken struct {
    bun.BaseModel `bun:"table:refresh_tokens,alias:at"`

    ID			int64      `bun:"id,pk,autoincrement"`
    HashToken	string     `bun:"hash_token,notnull,unique"`
    UserID		int64       `bun:"user_id,notnull"`
    ExpiresAt	time.Time  `bun:"expires_at,notnull"`
    UsedAt		*time.Time `bun:"used_at"`

    User *users.User `bun:"rel:belongs-to,join:user_id=id"`
}