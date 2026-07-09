package user

import (
	"golang.org/x/oauth2"
    "github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID 			int64 `bun:"id,pk,autoincrement"`
	Login 		string `bun:"login,unique,notnull"`
	Username 	string `bun:"name,unique,notnull"`
	Password 	string `bun:"password,notnull"`
	Email 		string `bun:"email,unique,notnull"`
	Is_verified bool `bun:"is_verified,notnull"`
	OAuth		*oauth2.Token `bun:"oauth,type:json"`
	Avatar 		*string `bun:"avatar"`
}
