package tool

import (
	"errors"
    "github.com/gin-gonic/gin"
)

func ConvertionUserId(c *gin.Context) (int64, error) {
	getId, ok := c.Get("userId")
	if !ok {
		return 0, errors.New("userId not found in context")
	}
	userId, ok := getId.(int64)
	if !ok {
		return 0, errors.New("userId not found in context")
	}
	return userId, nil
}