package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
)

/**
 * Middleware intercepts HTTP requests before they reach the handlers.
 * It is used to apply filters to specific endpoints—for example, checking
 * whether a user is authenticated via their token.
 *
 * Currently, this package serves as a placeholder and performs no action,
 * but it will gradually integrate these features in future commits.
 */

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			fmt.Printf("%s\n", err)
			c.AbortWithStatusJSON(401, gin.H{"state": "failed", "error_message": "unauthorized"})
			return
		}
		id, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"state": "failed", "error_message": "invalid token"})
			return
		}
		c.Set("userId", id)
		c.Next()
	}
}
