package handler

import (
	"os"
	"fmt"
	"time"
	"net/http"
	"backend/internal/tools"
	"github.com/uptrace/bun"
	"github.com/gin-gonic/gin"
	"backend/internal/repository"
)

func UserActivation(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		aToken, err := repository.UserActivation(c.Request.Context(), db, tokenString)
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

func RefreshToken(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(500, gin.H{"state": "failed",
                "error_message": "The user id doesn't exist"})
        }
		err := repository.RefreshToken(c.Request.Context(), db, userId.(int64))
		if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"state": "failed", "message": "You have to login"})
		} else {
            expiresAt := time.Now().Add(15 * time.Minute)
            newToken, tokenErr := tool.GenerateToken(userId.(int64), expiresAt, 
                []byte(os.Getenv("JWT_SECRET")))
            if tokenErr != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
                    "state": "failed", "message": "token creation went wrong"})
                return
            }
            c.JSON(http.StatusCreated, gin.H{
                "state": "success",
				"Authorization": newToken.Token})
		}
	}
}