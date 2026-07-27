package repository

import (
    "time"
    "errors"
	"context"
    "backend/internal/models"
	"github.com/uptrace/bun"
)

// ActivationToken
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

func ActivationToken(ctx context.Context, db *bun.DB,
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

func DeleteActivationToken(ctx context.Context, db *bun.DB,
    userId int64) error {
	_, err := db.NewDelete().
		Model((*model.ActivationToken)(nil)).
		Where("user_id = ?", userId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func SetUsedAt(ctx context.Context, db *bun.DB, userId int64) error {
    usedAt := time.Now()
    aToken := &model.ActivationToken{UsedAt: &usedAt}
    _, err := db.NewUpdate().
		Model(aToken).
		Column("used_at").
		Where("user_id = ?", userId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}



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
        return nil, errors.New("This token has expired, you have to login")
    } else if rToken.UsedAt != nil {
        return nil, errors.New("This token has been used already, you have to login")
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