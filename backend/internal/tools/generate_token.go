package tool

import (
	"backend/internal/models/users/tokens"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

func GenerateToken(userId int64, expiresAt time.Time, secretKey []byte) (*tokens.ActivationToken, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userId, 10),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "golang-api",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, tokenErr := token.SignedString(secretKey)
	if tokenErr != nil {
		return nil, tokenErr
	}
	aToken := &tokens.ActivationToken{UserID: userId,
		Token:     tokenString,
		ExpiresAt: expiresAt}
	return aToken, nil
}
