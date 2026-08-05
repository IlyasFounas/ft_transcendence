package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"context"

	"github.com/gin-gonic/gin"

	"backend/internal/models/users/tokens"
	"backend/internal/repositories"
	"backend/internal/tools"
)

type TokenHandler struct {
	userRepo  *repositories.UserRepository
	tokenRepo *repositories.TokenRepository
}

func NewTokenHandler(userRepo *repositories.UserRepository, tokenRepo *repositories.TokenRepository) *TokenHandler {
	return &TokenHandler{userRepo: userRepo, tokenRepo: tokenRepo}
}

// this function send an email to the client, to activate his account
func (h *TokenHandler) ActivationToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		aToken, err := h.tokenRepo.ActivationToken(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"state": "failed", "message": "account activation went wrong"})
			return
		}
		if err := h.tokenRepo.SetUsedAt(c.Request.Context(), aToken.UserID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"state": "failed", "message": "you already used this link"})
			return
		}
		if err := h.userRepo.SetIsVerified(c.Request.Context(), aToken.UserID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"state": "failed", "message": "user activation failed, contact the support"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"state": "success", "message": "your account has been activated"})
	}
}

func (h *TokenHandler) CreateRefreshToken(ctx context.Context, userId int64) (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	encodedKey := hex.EncodeToString(key)
	hashedKey := sha256.Sum256([]byte(encodedKey))
	hashToken := fmt.Sprintf("%x", hashedKey)
	rToken := &tokens.RefreshToken{
		HashToken: hashToken,
		UserID:    userId,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := h.tokenRepo.CreateRefreshToken(ctx, rToken); err != nil {
		return "", err
	}
	return encodedKey, nil
}

// we start a token in the db to refresh a short connection token that we give to the client
func (h *TokenHandler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		encodedKey := c.GetHeader("Refresh-Token")
		hashedKey := sha256.Sum256([]byte(encodedKey))
		hashToken := fmt.Sprintf("%x", hashedKey)
		rToken, err := h.tokenRepo.RefreshToken(c.Request.Context(), hashToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"state": "failed", "message": "You have to login"})
			return
		}
		if err := h.tokenRepo.SetRefreshUsedAt(c.Request.Context(), hashToken); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"state": "failed", "message": "You have to login"})
			return
		}
		expiresAt := time.Now().Add(15 * time.Minute)
		newToken, tokenErr := tool.GenerateToken(rToken.UserID, expiresAt,
			[]byte(os.Getenv("JWT_SECRET")))
		if tokenErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"state": "failed", "message": "token creation went wrong"})
			return
		}
		newEncodedKey, rTokenErr := h.CreateRefreshToken(c.Request.Context(), rToken.UserID)
		if rTokenErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"state": "failed", "message": "The refresh token creation failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"state":         "success",
			"Authorization": newToken.Token,
			"Refresh-Token": newEncodedKey,
		})
	}
}
