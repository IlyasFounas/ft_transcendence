package users

import (
	"github.com/uptrace/bun"
	"golang.org/x/oauth2"
	"time"
)

// this struct represent the Table user in the db
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID          int64         `bun:"id,pk,autoincrement"`
	Login       string        `bun:"login,unique,notnull"`
	Username    string        `bun:"name,unique,notnull"`
	Password    string        `bun:"password,notnull"`
	Email       string        `bun:"email,unique,notnull"`
	Is_verified bool          `bun:"is_verified,notnull"`
	OAuth       *oauth2.Token `bun:"oauth,type:json"`
	Avatar      *string       `bun:"avatar"`
	CreatedAt   time.Time     `bun:"created_at,default:current_timestamp"`
}
