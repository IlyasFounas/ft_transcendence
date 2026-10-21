package middleware

import (
	// "time"
	"context"
	"os"
	"fmt"
	"net/http"
	"log"
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
	tokenString, err := c.Cookie("Authorization")
    if err != nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you have to login"})
      return 
    }
    claims := jwt.MapClaims{}
    _, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // so we reject anything not signed with the HMAC family
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
          return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
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
	// exp, err := claims.GetExpirationTime()
	// if err != nil {
	//   fmt.Printf("%s\n", err)
    //   c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you have to login"})
    //   return
    // }
	// if exp == nil {
    //   fmt.Println("the claim expiration is malformed")
    //   c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you have to login"})
    //   return
	// }
	// if time.Now().After(exp.Time) {
    //   fmt.Println("the token has expired")
    //   c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you have to login"})
    //   return
    // }
    id, err := strconv.ParseInt(sub, 10, 64)
    if err != nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you have to login"})
      return
    }
    c.Set("userId", id)
    c.Next()
  }
}

// Tells whether a user is an admin
// UserRepository implements it, so the middleware package doesn't need to import the repositories.
type AdminChecker interface {
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

var adminChecker AdminChecker

// Way to check in the database that a user is an admin
// Called once in the main, right after the repositories are created
func SetAdminChecker(checker AdminChecker) {
	adminChecker = checker
}

// AdminMiddleWare only lets the admins through. It reads the userId set by AuthMiddleWare, so it must come after it:
// The admin flag is read from the database at every request (not from the token): taking the admin rights away
// from a user takes effect at once.
func AdminMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, ok := c.Get("userId")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you have to login"})
			return
		}
		userId, ok := value.(int64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you have to login"})
			return
		}
		if adminChecker == nil {
			// SetAdminChecker was never called: nobody gets through, rather than everybody
			log.Println("error : AdminMiddleWare is used but SetAdminChecker was never called")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		isAdmin, err := adminChecker.IsAdmin(c.Request.Context(), userId)
		if err != nil {
			log.Printf("error : %s\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "you are not allowed to do this"})
			return
		}
		c.Next()
	}
}