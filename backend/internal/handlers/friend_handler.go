package handlers

import (
	"net/http"
	"strconv"
	"errors"

	"github.com/gin-gonic/gin"
	"database/sql"

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

func (h *FriendHandler) AddFriend() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
				return
			}
		RequesterID := userId.(int64)
		username := c.PostForm("username")
		user, err := h.userRepo.UserByName(c.Request.Context(), username)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": "user not found"})
			return
		}
		if RequesterID == user.ID {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": "user sent bad request"})
			return
		}
		exists, err := h.friendRepo.FriendshipExists(c.Request.Context(), RequesterID, user.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound,
					gin.H{"message": "Friendship not found"})
				return
			}
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Failed to get friendship"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict,
				gin.H{"message": "Friendship is already existed"})
			return
		}
		friendObj := &friends.Friend{
			Requester_ID: RequesterID,
			Receiver_ID: user.ID,
		}
		_, err = h.friendRepo.AddFriend(c.Request.Context(), friendObj)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound,
					gin.H{"message": "user not found"})
				return
			}
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
				gin.H{"message": "No matching request"})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "accepted friend"})
	}
}

func (h *FriendHandler) RejectRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		rejectID := c.Param("id")
		rejectIDInt, err := strconv.ParseInt(rejectID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid ID"})
			return
		}
		updated, err := h.friendRepo.RejectRequest(c.Request.Context(), userId.(int64), rejectIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Failed to reject friend request"})
			return
		}
		if !updated {
			c.JSON(http.StatusNotFound,
				gin.H{"message": "No matching request"})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "rejected friend"})
	}
}

func (h *FriendHandler) CancelRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		cancelID := c.Param("id")
		cancelIDInt, err := strconv.ParseInt(cancelID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid ID"})
			return
		}
		updated, err := h.friendRepo.CancelRequest(c.Request.Context(), userId.(int64), cancelIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Failed to cancel friend request"})
			return
		}
		if !updated {
			c.JSON(http.StatusNotFound,
				gin.H{"message": "No matching request"})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "canceled friend request"})
	}
}

func (h *FriendHandler) RemoveFriend() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		removeID := c.Param("id")
		removeIDInt, err := strconv.ParseInt(removeID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid ID"})
			return
		}
		updated, err := h.friendRepo.RemoveFriend(c.Request.Context(), userId.(int64), removeIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Failed to remove friend"})
			return
		}
		if !updated {
			c.JSON(http.StatusNotFound,
				gin.H{"message": "No matching request"})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "removed friend"})
	}
}

// getter

func (h *FriendHandler) Friends() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		friendList, err := h.friendRepo.Friends(c.Request.Context(), userId.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Failed to get friend list"})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{
				"message": "success",
				"friends": friendList,
			})
	}
}

func (h *FriendHandler) FriendStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		username := c.Param("username")
		if username == "" {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": "username is required"})
			return
		}
		user, err := h.userRepo.UserByName(c.Request.Context(), username)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": "user not found"})
			return
		}
		searchUserID := user.ID
		friendStatus, err := h.friendRepo.FriendStatus(c.Request.Context(), userId.(int64), searchUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Failed to get friend status"})
			return
		}

		if friendStatus.Status == "none" {
            friendStatus.Name = user.Username
            friendStatus.Avatar = user.Avatar
			friendStatus.Rank = user.Rank
        }

		c.JSON(http.StatusOK,
			gin.H{
				"message": "success",
				"friendStatus": friendStatus,
			})
	}
}

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

func (h *FriendHandler) FriendPending() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		friendPending, err := h.friendRepo.FriendPending(c.Request.Context(), userId.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Failed to get friend pending list"})
			return 
		}
		c.JSON(http.StatusOK,
			gin.H{
				"message": "success",
				"friendPending": friendPending,
			})
	}
}
