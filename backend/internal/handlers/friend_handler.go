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

/**
 * func:
 * - AddFriend()
 * - RemoveFriend()
 * - AcceptRequest()
 * - RejectRequest()
 * - CancelRequest()
 * 
 * getter func:
 * - Friends(userID): get list of UserID friends
 * - IsFriend(userID): check if this id is friend
 * - FriendRequest()
 * - FriendPending()
 */


func (h *FriendHandler) AddFriend() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		username := c.PostForm("username")
		user, err := h.userRepo.UserByName(c.Request.Context(), username)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": "can't get the user by name"})
			return
		}
		friendObj := &friends.Friend{
			Requester_ID: userId.(int64),
			Receiver_ID: user.ID,
		}
		err = h.friendRepo.AddFriend(c.Request.Context(), friendObj)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "failed to add friend"})
			return
		}
		c.JSON(http.StatusCreated,
			gin.h{"message": "friend request sent"})
	}
}

