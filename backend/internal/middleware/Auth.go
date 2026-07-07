package auth

import (
    "github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
  return func(c *gin.Context) {
    // ... check token, session, etc.
    c.Next()
  }
}