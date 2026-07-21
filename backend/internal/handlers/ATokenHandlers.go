package handler

import (
	"net/http"
	"github.com/uptrace/bun"
	"github.com/gin-gonic/gin"
	"backend/internal/repository"
)

func UserActivation(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		aToken, err := repository.UserActivation(c.Request.Context(), db, tokenString)
		if err != nil {
			fmt.Printf("Token checking error : %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
            	"state": "failed", "message": "account activation went wrong"})
			return
		} else if aToken != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
            	"state": "failed", "message": "wrong link activation"})
			return 
		} else {
			err = repository.SetUsedAt(c.Request.Context(), db, aToken.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"state": "failed", "message": "wrong link activation"})
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