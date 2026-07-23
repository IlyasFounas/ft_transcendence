package repository

import (
    "time"
    "errors"
	"context"
    "backend/internal/models"
	"github.com/uptrace/bun"
)

// CreateActivationToken
func CreateAToken(ctx context.Context, db *bun.DB, aToken *model.ActivationToken,
    expiresAt time.Time) error {
	_, err := db.NewInsert().
			Model(aToken).
			Exec(ctx)
    if err != nil {
        return err
    }
	return nil
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

func CreateRefreshToken(ctx context.Context, db *bun.DB, rToken *model.RefreshToken) error {
	_, err := db.NewInsert().
				Model(rToken).
				Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func RefreshToken(ctx context.Context, db *bun.DB, userId int64) error {
    rToken := new(model.RefreshToken)
	err := db.NewSelect().
            Model(rToken).
            Column("id", "user_id", "expires_at", "used_at").
            Where("user_id = ?", userId).
			Scan(ctx)
	if err != nil {
		return err
	}
	return nil
}