package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"backend/internal/models/users"
	"backend/internal/repositories"
	"backend/internal/tools/encryption"
	"backend/internal/tools/mailing"
	"backend/internal/tools/token_generation"
)

/**
 * README:
 * Handlers make the bridge between the router and the repository functions.
 * It's important to split different features into different files to improve visibility
 * and maintainability.
 *
 * Basically, the handlers call the repository functions, manage any errors,
 * and return the result to the frontend (in JSON, plain text, etc.).
 */
type UserHandler struct {
	userRepo     *repositories.UserRepository
	tokenRepo    *repositories.TokenRepository
	tokenHandler *TokenHandler
}

func NewUserHandler(userRepo *repositories.UserRepository, tokenRepo *repositories.TokenRepository, tokenHandler *TokenHandler) *UserHandler {
	return &UserHandler{userRepo: userRepo, tokenRepo: tokenRepo, tokenHandler: tokenHandler}
}

func (h *UserHandler) createActivationToken(c *gin.Context, expiresAt time.Time,
	userId int64) error {
	tokenObj, tokenErr := token_generation.GenerateToken(userId, expiresAt,
		[]byte(os.Getenv("ACTIVATION_JWT_SECRET")))
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"state": "failed", "message": "token creation went wrong"})
		return tokenErr
	}
	err := h.tokenRepo.DeleteActivationToken(c.Request.Context(), userId)
	if err != nil {
		fmt.Printf("Token creation error message : %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"state": "failed", "message": "token creation went wrong"})
		return err
	}
	err = h.tokenRepo.CreateAToken(c.Request.Context(), tokenObj, expiresAt)
	if err != nil {
		fmt.Printf("Token creation error message : %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"state": "failed", "message": "token creation went wrong"})
		return err
	}
	if err := mailing.SendActivationEmail(c.PostForm("email"), tokenObj.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"state": "failed", "message": "link activation sending went wrong"})
		return err
	}
	return nil
}

func (h *UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		hashpassword, hashErr := encryption.HashPassword(c.PostForm("password"))
		if hashErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"state": "failed", "message": "password encryption failed"})
			return
		}
		userObj := &users.User{Login: c.PostForm("login"),
			Password:    hashpassword,
			Username:    c.PostForm("username"),
			Email:       c.PostForm("email"),
			Is_verified: false}

		user, err := h.userRepo.CreateUser(c.Request.Context(), userObj)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"state": "failed", "message": "database exection failed"})
			return
		} else {
			expiresAt := time.Now().Add(24 * time.Hour)
			if err := h.createActivationToken(c, expiresAt, user.ID); err != nil {
				return
			}
			c.JSON(http.StatusCreated, gin.H{
				"state": "success", "message": "you just receive an activation link on your email"})
		}
	}
}

func (h *UserHandler) LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		state, userId, err := h.userRepo.LoginUser(c.Request.Context(),
			c.PostForm("login"), c.PostForm("password"))
		if err != nil && userId != 0 {
			expiresAt := time.Now().Add(24 * time.Hour)
			if err := h.createActivationToken(c, expiresAt, userId); err != nil {
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"state": "failed",
				"message": "an verification mail has been sent"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
			return
		}
		if !state {
			c.JSON(http.StatusOK, gin.H{
				"state":   "failed",
				"message": "The login or the password is wrong"})
		} else {
			// generate the jwt token here
			expiresAt := time.Now().Add(15 * time.Minute)
			tokenObj, tokenErr := token_generation.GenerateToken(userId, expiresAt,
				[]byte(os.Getenv("JWT_SECRET")))
			if tokenErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"state": "failed", "message": "token creation went wrong"})
				return
			}
			encodedKey, rTokenErr := h.tokenHandler.CreateRefreshToken(c.Request.Context(), userId)
			if rTokenErr != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"state": "failed", "message": "The refresh token creation failed"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"state":         "success",
				"Authorization": tokenObj.Token,
				"Refresh-Token": encodedKey})
		}
	}
}

// Getters
func (h *UserHandler) GetUser() gin.HandlerFunc {
	//return almost all user infos
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"state": "failed",
				"error_message": "The user id doesn't exist"})
		}
		userData, err := h.userRepo.GetUser(c.Request.Context(), userId.(int64))
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"state": "failed", "error_message": err})
		}
		if userData == nil {
			c.JSON(http.StatusNotFound, gin.H{"state": "failed", "message": "The user doesn't exist"})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"state": "success",
				"user":  userData,
			})
		}
	}
}
func (h *UserHandler) GetName() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"state": "failed",
				"error_message": "The user id doesn't exist"})
		}
		name, err := h.userRepo.GetName(c.Request.Context(), userId.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
		}
		if name == nil {
			c.JSON(http.StatusNotFound, gin.H{"state": "failed"})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"state": "success",
				"name":  name})
		}
	}
}
func (h *UserHandler) GetLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"state": "failed",
				"error_message": "The user id doesn't exist"})
		}
		login, err := h.userRepo.GetLogin(c.Request.Context(), userId.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
		}
		if login == nil {
			c.JSON(http.StatusNotFound, gin.H{"stsate": "failed"})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"state": "success",
				"login": login})
		}
	}
}
func (h *UserHandler) GetEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"state": "failed",
				"error_message": "The user id doesn't exist"})
		}
		email, err := h.userRepo.GetEmail(c.Request.Context(), userId.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
		}
		if email == nil {
			c.JSON(http.StatusNotFound, gin.H{"state": "failed", "message": "the user have no email"})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"state": "success",
				"email": email})
		}
	}
}
func (h *UserHandler) GetOAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"state": "failed",
				"error_message": "The user id doesn't exist"})
		}
		oauth, err := h.userRepo.GetOAuth(c.Request.Context(), userId.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
		}
		if oauth == nil {
			c.JSON(http.StatusNotFound, gin.H{"state": "failed", "message": "oauth not found"})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"state": "success",
				"OAuth": oauth})
		}
	}
}

// Setters
func (h *UserHandler) SetName() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"state": "failed",
				"error_message": "The user id doesn't exist"})
		}
		state, err := h.userRepo.SetName(c.Request.Context(), userId.(int64), c.PostForm("name"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
		}
		if !state {
			c.JSON(http.StatusOK, gin.H{
				"state": "failed"})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"state": "success"})
		}
	}
}
func (h *UserHandler) SetLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"state": "failed",
				"error_message": "The user id doesn't exist"})
		}
		state, err := h.userRepo.SetLogin(c.Request.Context(), userId.(int64), c.PostForm("login"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
		}
		if !state {
			c.JSON(http.StatusOK, gin.H{
				"state": "failed"})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"state": "success"})
		}
	}
}
func (h *UserHandler) SetAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"state": "failed",
				"error_message": "The user id doesn't exist"})
		}
		state, err := h.userRepo.SetAvatar(c.Request.Context(), userId.(int64), c.PostForm("avatar"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
		}
		if !state {
			c.JSON(http.StatusOK, gin.H{
				"state": "failed"})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"state": "success"})
		}
	}
}
