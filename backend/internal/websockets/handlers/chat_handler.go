package ws_handlers

import (
	"backend/internal/models/users"
	"backend/internal/models/chats"
	"backend/internal/models/friends"
	"backend/internal/repositories"
	"context"
	"log"
	"strings"
	"time"
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
	userId int64, targetId int64, message string) (*chats.Chat, bool) {
	if strings.TrimSpace(message) == "" {
		log.Printf("error : empty message")
		return nil, false
	}
	if len([]rune(message)) > 1000 {
		log.Printf("error : message is too long")
		return nil, false
	}
	target, err := h.userRepo.User(ctx, targetId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, false
	}
	if target == nil {
		log.Printf("error : destination user id is invalid")
		return nil, false
	}
	friendStatus, err := h.friendRepo.FriendStatus(ctx, userId, targetId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, false
	}
	if friendStatus.Status != "friends" {
		log.Println("error : the friendship doesn't exist")
		return nil, false
	}
	chatObj := &chats.Chat{
		SRC_user_ID: userId,
		DST_user_ID: targetId,
		Message: message,
		CreatedAt: time.Now(),
	}
	_, err = h.chatRepo.Send(ctx, chatObj)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, false
	}
	return chatObj, true
}

func (h *ChatHandler) Receive(ctx context.Context,
	userId int64, friendId int64, actualMsgNb int64) ([]*chats.Chat, bool) {
	friend, err := h.userRepo.User(ctx, friendId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, false
	}
	if friend == nil {
		log.Printf("error : destination user id is invalid")
		return nil, false
	}
	friendStatus, err := h.friendRepo.FriendStatus(ctx, userId, friendId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, false
	}
	if friendStatus == nil {
		log.Println("error : friendStatus is NULL")
		return nil, false
	}
	if friendStatus.Status != "friends" {
		log.Println("error : the friendship doesn't exist")
		return nil, false
	}
	// Receive 20 older messages
	chats, err := h.chatRepo.Receive(ctx, userId, friendId, actualMsgNb)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, false
	}
	return chats, true
}

func (h *ChatHandler) FriendList(ctx context.Context,
	userId int64) ([]*friends.FriendResponse, bool) {
	friends, err := h.friendRepo.Friends(ctx, userId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, false
	}
	return friends, true
}

func (h *ChatHandler) GetUser(ctx context.Context,
	userId int64) (*users.User, bool) {
	user, err := h.userRepo.User(ctx, userId);
	if err != nil {
		log.Printf("error : %s\n", err);
		return nil, false
	}
	if user == nil {
		log.Printf("User doesn't exist\n")
		return nil, false
	}
	return user, true
}