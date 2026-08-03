package handler

import (
	"os"
	"fmt"
	"time"
	"context"
	"net/http"
	// "strconv"
    "backend/internal/models"
	"backend/internal/tools"
	"github.com/uptrace/bun"
	"github.com/gin-gonic/gin"
	"backend/internal/repository"
    // "github.com/golang-jwt/jwt/v5"
	"crypto/rand"
	"encoding/hex"
	"crypto/sha256"
)

// this function send an email to the client, to activate his account
func ActivationToken(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		aToken, err := repository.ActivationToken(c.Request.Context(), db, tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"state": "failed", "message": "account activation went wrong"})
			return
		} else {
			fmt.Printf("%d\n", aToken.UserID)
			err = repository.SetUsedAt(c.Request.Context(), db, aToken.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"state": "failed", "message": "you already used this link"})
				return
			}
			err = repository.SetIsVerified(c.Request.Context(), db, aToken.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"state": "failed", "message": "user activation failed, contact the support"})
				return 
			}
            c.JSON(http.StatusCreated, gin.H{
                "state": "success", "message": "your account has been activated"})
		}
	}
}

func CreateRefreshToken(ctx context.Context, db *bun.DB,
        userId int64) (string, error) {
	key := make([]byte, 32)
	rand.Read(key)
    encodedKey := hex.EncodeToString(key)
    hashedKey := sha256.Sum256([]byte(encodedKey))
	hashToken := fmt.Sprintf("%x", hashedKey)
	rToken := &model.RefreshToken{
		HashToken: hashToken,
		UserID: userId,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err := repository.CreateRefreshToken(ctx, db, rToken)
	if err != nil {
		return "", err
	}
	return encodedKey, nil
}

// we start a token in the db, to refresh a short connection token that we give to the client
func RefreshToken(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
    	encodedKey := c.GetHeader("Refresh-Token")
		hashedKey := sha256.Sum256([]byte(encodedKey))
		hashToken := fmt.Sprintf("%x", hashedKey)
		rToken, err := repository.RefreshToken(c.Request.Context(), db, hashToken)
		if err != nil {
			fmt.Printf("\n\n%s\n\n\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"state": "failed", "message": "You have to login"})
			return 
		} else {
			err = repository.SetRefreshUsedAt(c.Request.Context(), db, hashToken)
			if err != nil {
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
            encodedKey, rTokenErr := CreateRefreshToken(c.Request.Context(), db, rToken.UserID)
            if rTokenErr != nil {
                c.JSON(http.StatusUnauthorized, gin.H{
				    "state": "failed", "message": "The refresh token creation failed"})
                return
            }
            c.JSON(http.StatusOK, gin.H{
                "state": "success",
				"Authorization": newToken.Token,
                "Refresh-Token": encodedKey})
		}
	}
}