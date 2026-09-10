package handlers

import (
	"net/http"
	"strconv"

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
 * - FriendRequests()
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
		_, err = h.friendRepo.AddFriend(c.Request.Context(), friendObj)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Failed to add friend"})
			return
		}
		c.JSON(http.StatusCreated,
			gin.H{"message": "friend request sent"})
	}
}

func (h *FriendHandler) AcceptRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		requesterID := c.Param("id")
		requesterIDInt, err := strconv.ParseInt(requesterID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid requester ID"})
			return
		}
		updated, err := h.friendRepo.AcceptRequest(c.Request.Context(), userId.(int64), requesterIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Failed to accept friend request"})
			return
		}
		if !updated {
			c.JSON(http.StatusNotFound,
				gin.H{"message": "No matching pending request"})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "accepted friend"})
	}
}

// getter

// func (h *FriendHandler) Friends() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userId, ok := c.Get("userId")
// 		if !ok {
// 			c.JSON(http.StatusUnauthorized,
// 				gin.H{"message": "The userId is empty"})
// 			return
// 		}
// 		friends, err := h.friendRepo.Friends(c.Request.Context(), userId.(int64))
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError,
// 				gin.H{"message": "Failed to get friend list"})
// 			return
// 		}
// 		c.JSON(http.StatusOK,
// 			gin.H{
// 				"message": "success",
// 				"friends": friends,
// 			})
// 	}
// }

func (h *FriendHandler) FriendRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		friendRequests, err := h.friendRepo.FriendRequests(c.Request.Context(), userId.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Failed to get friend requests"})
			return 
		}
		c.JSON(http.StatusOK,
			gin.H{
				"message": "success",
				"friendRequests": friendRequests,
			})
	}
}