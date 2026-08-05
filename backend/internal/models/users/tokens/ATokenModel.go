package tokens

import (
	"time"
    "github.com/uptrace/bun"
    "backend/internal/models/users"
)

type ActivationToken struct {
    bun.BaseModel `bun:"table:activation_tokens,alias:at"`

    ID        int64      `bun:"id,pk,autoincrement"`
    UserID    int64      `bun:"user_id,notnull"`
    Token     string     `bun:"token,notnull,unique"`
    ExpiresAt time.Time  `bun:"expires_at,notnull"`
    UsedAt    *time.Time `bun:"used_at"`

    // we allow joins between activation_tokens and users tables
    User *users.User `bun:"rel:belongs-to,join:user_id=id"`
}