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