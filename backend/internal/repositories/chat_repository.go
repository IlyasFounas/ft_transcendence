package repositories

import (
	"backend/internal/models/chats"
	"context"

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
	var chats []*chats.Chat
	err := r.db.NewSelect().
		Model(&chats).
		Order("id DESC").
		Where("id >= ? ", actualMsgNb - 20).
		Where("id < ? ", actualMsgNb).
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.WhereOr("src_user_id = ? && dst_user_id = ?", userId, friendId).
				WhereOr("src_user_id = ? AND dst_user_id = ?", friendId, userId)
		}).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return chats, nil
}