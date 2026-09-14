package handlers

import (
	"log"
	"context"
	"backend/internal/repositories"
)

type ChatHandler struct {
	chatRepo	*repositories.ChatRepository
	userRepo	*repositories.UserRepository
}

func NewChatHandler (chatRepo *repositories.ChatRepository,
	userRepo	*repositories.UserRepository) *ChatHandler {
	return &ChatHandler{chatRepo: chatRepo,
		userRepo: userRepo}
}