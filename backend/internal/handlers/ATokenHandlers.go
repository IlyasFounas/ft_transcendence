package handler

import (
	"fmt"
	"net/http"
	"github.com/uptrace/bun"
	"github.com/gin-gonic/gin"
	"backend/internal/repository"
)

func UserActivation(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		fmt.Printf("\n\nToken ===== : %s\n\n", tokenString)
		aToken, err := repository.UserActivation(c.Request.Context(), db, tokenString)
		fmt.Printf("\n\nToken checking error : %s\n\n", err)
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