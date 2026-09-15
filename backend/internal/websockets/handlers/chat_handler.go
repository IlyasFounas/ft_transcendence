package ws_handlers

import (
	"log"
	"context"
	"strings"
	"backend/internal/models/chats"
	"backend/internal/repositories"
)

type ChatHandler struct {
	chatRepo	*repositories.ChatRepository
	friendRepo	*repositories.FriendRepository
	userRepo	*repositories.UserRepository
}

func NewChatHandler (chatRepo *repositories.ChatRepository,
	friendRepo	*repositories.FriendRepository,
	userRepo	*repositories.UserRepository) *ChatHandler {
	return &ChatHandler{chatRepo: chatRepo,
		friendRepo: friendRepo,
		userRepo: userRepo}
}

func (h *ChatHandler) Send(ctx context.Context, 
	userId int64, target_id int64, message string) bool {
	if strings.TrimSpace(message) == "" {
		log.Printf("error : empty message")
		return false
	}
	if len([]rune(message)) > 1000 {
		log.Printf("error : message is too long")
		return false
	}
	target, err := h.userRepo.User(ctx, target_id)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	if target == nil {
		log.Printf("error : destination user id is invalid")
		return false
	}
	chatObj := &chats.Chat{
		SRC_user_ID: userId,
		DST_user_ID: target_id,
		Message: message,
	}
	_, err = h.chatRepo.Send(ctx, chatObj)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	return true
} 