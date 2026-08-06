package repositories

import (
	"backend/internal/models/users/tokens"
	"context"
	"errors"
	"github.com/uptrace/bun"
	"time"
)

type RefreshTokenRepository struct {
	db *bun.DB
}

func NewRefreshTokenRepository(db *bun.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// RefreshToken
func (tr *RefreshTokenRepository) CreateRefreshToken(ctx context.Context, rToken *tokens.RefreshToken) error {
	_, err := tr.db.NewInsert().
		Model(rToken).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (tr *RefreshTokenRepository) RefreshToken(ctx context.Context,
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

func (tr *RefreshTokenRepository) SetRefreshUsedAt(ctx context.Context, hashToken string) error {
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
