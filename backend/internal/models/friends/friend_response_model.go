package friends

import "time"

type FriendResponse struct {
	ID			int64   	`json:"id"`
	Name		string    	`json:"name"`
	Avatar    	*string   	`json:"avatar"`
	CreatedAt 	time.Time 	`json:"created_at"`
}