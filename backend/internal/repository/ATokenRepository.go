package repository

import (
    "os"
    "time"
    "errors"
	"context"
    "strconv"
    "backend/internal/models"
	"github.com/uptrace/bun"
    "github.com/golang-jwt/jwt/v5"
)

// CreateActivationToken
func CreateAToken(ctx context.Context, user *model.User, db *bun.DB) (*string, error) {
    expiresAt := time.Now().Add(24 * time.Hour)
    claims := jwt.RegisteredClaims{
            Subject: strconv.FormatInt(user.ID, 10),
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "golang-api",
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, tokenErr := token.SignedString([]byte(os.Getenv("ACTIVATION_JWT_SECRET")))
    if tokenErr != nil {
        return nil, tokenErr 
    }
    aToken := &model.ActivationToken {UserID: user.ID,
        Token: tokenString,
        ExpiresAt: expiresAt }
	_, err := db.NewInsert().
			Model(aToken).
			Exec(ctx)
    if err != nil {
        return nil, err
    }
	return &tokenString, nil
}

func UserActivation(ctx context.Context, db *bun.DB,
    tokenString string) (*model.ActivationToken, error) {
    aToken := new(model.ActivationToken)
	err := db.NewSelect().
            Model(aToken).
            Column("id", "user_id", "expires_at", "used_at").
            Where("token = ?", tokenString).
			Scan(ctx)
    if err != nil {
        return nil, err
    } else if aToken.UsedAt != nil {
        return nil, errors.New("This token is already used")
    } else if aToken.ExpiresAt.Before(time.Now()) {
        return nil, errors.New("This token has expired")
    }
	return aToken, nil
}

func SetUsedAt(ctx context.Context, db *bun.DB, userId int64) error {
    usedAt := time.Now()
    aToken := &model.ActivationToken{UsedAt: &usedAt}
    _, err := db.NewUpdate().
		Model(aToken).
		Column("used_at").
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}