package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/models/friends"
	"backend/internal/repositories"
)

type FriendHandler struct {
	friendRepo	*repositories.FriendRepository
	userRepo	*repositories.UserRepository
}

func NewFriendHandler (friendRepo *repositories.FriendRepository,
	userRepo	*repositories.UserRepository) *FriendHandler {
	return &FriendHandler{friendRepo: friendRepo,
		userRepo: userRepo}
}