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