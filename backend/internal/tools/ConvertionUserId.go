package tool

import (
	"errors"
	"strconv"
)

func ConvertionUserId(userIdString string) (int64, error) {
    userId, err := strconv.ParseInt(userIdString, 10, 64)
    if err != nil {
        return 0, errors.New("conversion failed")
    }
	return userId, nil
}