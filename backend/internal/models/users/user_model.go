package users

import (
	"github.com/uptrace/bun"
	"time"
	"backend/internal/models/games"
)

// this struct represent the Table user in the db
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID          int64         	`bun:"id,pk,autoincrement"`
	Login       *string        	`bun:"login,unique"`
	Username    string        	`bun:"name,unique,notnull"`
	GoogleId 	*string 	  	`bun:"google_id,unique"`
	Password 	*string 	  	`bun:"password"`
	Email       string        	`bun:"email,unique,notnull"`
	Is_verified bool          	`bun:"is_verified,notnull"`
	Avatar      *string       	`bun:"avatar,default:'/avatars/avatar1.png'"`
	XP   		int64  			`bun:"xp,notnull,default:0"`
    Rank 		string 			`bun:"rank,notnull,default:'Junior Agent'"`
	CreatedAt   time.Time     	`bun:"created_at,default:current_timestamp"`

	Lobby 		[]*games.Lobby	`bun:"rel:has-many,join:id=user_id"`
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}