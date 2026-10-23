package handlers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"backend/internal/models/users/tokens"
	"backend/internal/repositories"
	"backend/internal/tools/cookie"
	"backend/internal/tools/token_generation"
)

type TokenHandler struct {
	userRepo   *repositories.UserRepository
	aTokenRepo *repositories.ActivationTokenRepository
	rTokenRepo *repositories.RefreshTokenRepository
}

func NewTokenHandler(userRepo *repositories.UserRepository,
	aTokenRepo *repositories.ActivationTokenRepository,
	rTokenRepo *repositories.RefreshTokenRepository) *TokenHandler {
	return &TokenHandler{userRepo: userRepo,
		aTokenRepo: aTokenRepo,
		rTokenRepo: rTokenRepo}
}

// this function send an email to the client, to activate his account
func (h *TokenHandler) ActivationToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		aToken, err := h.aTokenRepo.ActivationToken(c.Request.Context(), tokenString)
		if err != nil {
			log.Printf("error : %s\nerror : account activation went wrong\n", err)
			c.Redirect(http.StatusPermanentRedirect, os.Getenv("PUBLIC_URL"))
			return
		}
		if err := h.aTokenRepo.SetUsedAt(c.Request.Context(), aToken.UserID); err != nil {
			log.Printf("this link is already used\n")
			c.Redirect(http.StatusPermanentRedirect, os.Getenv("PUBLIC_URL"))
			return
		}
		if err := h.userRepo.SetIsVerified(c.Request.Context(), aToken.UserID); err != nil {
			log.Printf("user activation failed, contact the support\n")
			c.Redirect(http.StatusPermanentRedirect, os.Getenv("PUBLIC_URL"))
			return
		}
		loginPage := fmt.Sprintf("%s/login", os.Getenv("PUBLIC_URL"))
		c.Redirect(http.StatusPermanentRedirect, loginPage)
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
	if err := h.rTokenRepo.CreateRefreshToken(ctx, rToken); err != nil {
		return "", err
	}
	return encodedKey, nil
}

func (h *TokenHandler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		encodedKey, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You have to login"})
			return
		}
		hashedKey := sha256.Sum256([]byte(encodedKey))
		hashToken := fmt.Sprintf("%x", hashedKey)
		cookie.Clear(c, "/api/auth", "refresh_token")
		// check the refresh token
		rToken, err := h.rTokenRepo.RefreshToken(c.Request.Context(), hashToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You have to login"})
			return
		}
		if rToken.UsedAt != nil {
			err := h.rTokenRepo.DestroyTokens(c.Request.Context(), rToken.UserID)
			if err != nil {
				c.JSON(http.StatusUnauthorized,
					gin.H{"message": "Contact the support"})
				return
			}
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "This refresh token has been used, you have to login"})
			return
		}
		if err := h.rTokenRepo.SetRefreshUsedAt(c.Request.Context(), hashToken); err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "You have to login"})
			return
		}
		// creating the refresh_token
		newEncodedKey, rTokenErr := h.CreateRefreshToken(c.Request.Context(), rToken.UserID)
		if rTokenErr != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "You have to login"})
			return
		}
		cookie.Set(c, "/api/auth", "refresh_token",
			newEncodedKey, 7*24*time.Hour, true)
		// creating the access_token
		expiresAt := time.Now().Add(15 * time.Minute)
		newToken, tokenErr := token_generation.GenerateToken(rToken.UserID, expiresAt,
			[]byte(os.Getenv("JWT_SECRET")))
		if tokenErr != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "You have to login"})
			return
		}
		cookie.Set(c, "/api", "Authorization",
			newToken.Token, 15*time.Minute, true)
		c.JSON(http.StatusOK, gin.H{
			"message": "success"})
	}
}
