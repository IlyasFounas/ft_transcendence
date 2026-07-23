package handler

import (
	"os"
	"fmt"
	"time"
	"net/http"
	"strconv"
	"backend/internal/tools"
	"github.com/uptrace/bun"
	"github.com/gin-gonic/gin"
	"backend/internal/repository"
    "github.com/golang-jwt/jwt/v5"
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
    	tokenString := c.GetHeader("Authorization")
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil })
		userId, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"state": "failed", "error_message": "invalid token"})
			return
		}
		fmt.Printf("\n\n\n\n\n%d\n\n\n\n\n", userId)
		err = repository.RefreshToken(c.Request.Context(), db, userId)
		if err != nil {
			fmt.Printf("\n\n%s\n\n\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"state": "failed", "message": "You have to login"})
		} else {
            expiresAt := time.Now().Add(15 * time.Minute)
            newToken, tokenErr := tool.GenerateToken(userId, expiresAt, 
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