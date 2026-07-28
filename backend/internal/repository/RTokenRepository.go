package repository

import (
    "time"
    "errors"
	"context"
    "backend/internal/models"
	"github.com/uptrace/bun"
)

// RefreshToken
func CreateRefreshToken(ctx context.Context, db *bun.DB, rToken *model.RefreshToken) error {
	_, err := db.NewInsert().
		Model(rToken).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func RefreshToken(ctx context.Context, db *bun.DB,
		hashToken string) (*model.RefreshToken, error) {
    rToken := new(model.RefreshToken)
	err := db.NewSelect().
            Model(rToken).
            Column("id", "user_id", "expires_at", "used_at").
            Where("hash_token = ?", hashToken).
			Scan(ctx)
	if err != nil {
		return nil, err
	} else if rToken.ExpiresAt.Before(time.Now()) {
        return nil, errors.New("You have to login")
    }
	return rToken, nil
}

func SetRefreshUsedAt(ctx context.Context, db *bun.DB, hashToken string) error {
    usedAt := time.Now()
    rToken := &model.RefreshToken{UsedAt: &usedAt}
    _, err := db.NewUpdate().
		Model(rToken).
		Column("used_at").
		Where("hash_token = ?", hashToken).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}