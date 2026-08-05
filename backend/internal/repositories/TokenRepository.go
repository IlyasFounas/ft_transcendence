package repositories

import (
	"context"
    "errors"
    "time"
    "backend/internal/models/users/tokens"
	"github.com/uptrace/bun"
)

type TokenRepository struct {
	db         *bun.DB
}

func NewTokenRepository(db *bun.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// ActivationToken
func (tr *TokenRepository) CreateAToken(ctx context.Context, aToken *tokens.ActivationToken,
    expiresAt time.Time) error {
	_, err := tr.db.NewInsert().
			Model(aToken).
			Exec(ctx)
    if err != nil {
        return err
    }
	return nil
}

func (tr *TokenRepository) ActivationToken(ctx context.Context,
    tokenString string) (*tokens.ActivationToken, error) {
    aToken := new(tokens.ActivationToken)
	err := tr.db.NewSelect().
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

func (tr *TokenRepository) DeleteActivationToken(ctx context.Context,
    userId int64) error {
	_, err := tr.db.NewDelete().
		Model((*tokens.ActivationToken)(nil)).
		Where("user_id = ?", userId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TokenRepository) SetUsedAt(ctx context.Context, userId int64) error {
    usedAt := time.Now()
    aToken := &tokens.ActivationToken{UsedAt: &usedAt}
    _, err := tr.db.NewUpdate().
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
func (tr *TokenRepository) CreateRefreshToken(ctx context.Context, rToken *tokens.RefreshToken) error {
	_, err := tr.db.NewInsert().
		Model(rToken).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TokenRepository) RefreshToken(ctx context.Context,
		hashToken string) (*tokens.RefreshToken, error) {
    rToken := new(tokens.RefreshToken)
	err := tr.db.NewSelect().
            Model(rToken).
            Column("id", "user_id", "expires_at", "used_at").
            Where("hash_token = ?", hashToken).
			Scan(ctx)
	if err != nil {
		return nil, err
	} else if rToken.ExpiresAt.Before(time.Now()) {
        return nil, errors.New("This token has expired, you have to login")
    }
	return rToken, nil
}

func (tr *TokenRepository) SetRefreshUsedAt(ctx context.Context, hashToken string) error {
    usedAt := time.Now()
    rToken := &tokens.RefreshToken{UsedAt: &usedAt}
    _, err := tr.db.NewUpdate().
		Model(rToken).
		Column("used_at").
		Where("hash_token = ?", hashToken).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}