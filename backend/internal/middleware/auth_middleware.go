package middleware

import (
    "os"
    "fmt"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
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
    // tokenString := c.GetHeader("Authorization")
		tokenString, err := c.Cookie("Authorization")
    if err != nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you have to login"})
      return 
    }
    claims := jwt.MapClaims{}
    _, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil })
    if err != nil {
      fmt.Printf("%s\n", err)
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you have to login"})
      return
    }
    // avoids panic error when the token is valid but the "sub" claim is missing or not a string
    sub, ok := claims["sub"].(string)
    if !ok {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you have to login"})
      return
    }
    id, err := strconv.ParseInt(sub, 10, 64)
    if err != nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you have to login"})
      return
    }
    c.Set("userId", id)
    c.Next()
  }
}