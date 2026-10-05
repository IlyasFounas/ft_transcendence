package repositories

import (
	"backend/internal/models/chats"
	"context"
	"database/sql"
	"errors"
	"time"
	"github.com/uptrace/bun"
)

type ChatRepository struct {
	db *bun.DB
}

func NewChatRepository(db *bun.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Send(ctx context.Context, chat *chats.Chat) (*chats.Chat, error) {
	_, err := r.db.NewInsert().
		Model(chat).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (r *ChatRepository) Receive(ctx context.Context, userId int64,
	friendId int64, actualMsgNb int64) ([]*chats.Chat, error) {
	var chatsMessage []*chats.Chat
	request := r.db.NewSelect().
		Model(&chatsMessage).
		Order("id DESC").
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.WhereOr("src_user_id = ? AND dst_user_id = ?", userId, friendId).
				WhereOr("src_user_id = ? AND dst_user_id = ?", friendId, userId)
		}).
		Limit(20)
	if actualMsgNb > 0 {
		request = request.Where("id < ?", actualMsgNb)
	}
	err := request.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return chatsMessage, nil
}

func (r *ChatRepository) LastFriendMessage(ctx context.Context, 
	userId, friendId int64) (string, time.Time, error) {
	var date time.Time
	var lastMessage string
	err := r.db.NewSelect().
		Model((*chats.Chat)(nil)).
		Column("message", "created_at").
		Where("src_user_id = ? AND dst_user_id = ?", friendId, userId).
		Order("id DESC").
		Limit(1).
		Scan(ctx, &lastMessage, &date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", time.Now(), nil
		}
		return "", time.Now(), err
	}
	return lastMessage, date, nil
}