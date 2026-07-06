package user

import (
	"golang.org/x/oauth2"
    "github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID 			int64 `bun:"id,pk,autoincrement"`
	Name 		string `bun:"name,unique,notnull"`
	Password 	string `bun:"password,notnull"`
	Email 		string `bun:"email,unique,notnull"`
	OAuth		*oauth2.Token `bun:"oauth,type:json"`
	Avatar 		*string `bun:"avatar"`
}
