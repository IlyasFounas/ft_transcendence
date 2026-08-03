package tool

import (
	"time"
	"strconv"
    "backend/internal/models"
    "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId int64, expiresAt time.Time, secretKey []byte) (*model.ActivationToken, error) {
	claims := jwt.RegisteredClaims{
            Subject: strconv.FormatInt(userId, 10),
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "golang-api",
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, tokenErr := token.SignedString(secretKey)
    if tokenErr != nil {
        return nil, tokenErr 
    }
    aToken := &model.ActivationToken {UserID: userId,
    Token: tokenString,
    ExpiresAt: expiresAt }
	return aToken, nil
}