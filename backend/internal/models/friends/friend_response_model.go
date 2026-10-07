package friends

import "time"

type FriendResponse struct {
	ID			int64   	`json:"id"`
	UserID     	int64     	`json:"user_id"`
	Name		string    	`json:"name"`
	Avatar    	*string   	`json:"avatar"`
	Rank		string		`json:"rank"`
	CreatedAt 	time.Time 	`json:"created_at"`
}

type FriendStatusResponse struct {
    UserID      int64   `json:"user_id"`
    Name        string  `json:"name"`
    Avatar      *string `json:"avatar"`
	Rank		string	`json:"rank"`
    Status      string  `json:"status"`
    RequesterID int64   `json:"-"`
}