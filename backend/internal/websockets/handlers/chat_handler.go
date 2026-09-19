package ws_handlers

import (
	"backend/internal/models/chats"
	"backend/internal/models/friends"
	"backend/internal/repositories"
	"context"
	"log"
	"strings"
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
	userId int64, targetId int64, message string) bool {
	if strings.TrimSpace(message) == "" {
		log.Printf("error : empty message")
		return false
	}
	if len([]rune(message)) > 1000 {
		log.Printf("error : message is too long")
		return false
	}
	target, err := h.userRepo.User(ctx, targetId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	if target == nil {
		log.Printf("error : destination user id is invalid")
		return false
	}
	friendStatus, err := h.friendRepo.FriendStatus(ctx, userId, targetId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	if friendStatus.Status != "friends" {
		log.Println("error : the friendship doesn't exist")
		return false
	}
	chatObj := &chats.Chat{
		SRC_user_ID: userId,
		DST_user_ID: targetId,
		Message: message,
	}
	_, err = h.chatRepo.Send(ctx, chatObj)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	return true
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
	// we want the 20 more messages
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