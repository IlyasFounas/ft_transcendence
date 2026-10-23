package repositories

import (
	"backend/internal/models/users/tokens"
	"context"
	"errors"
	"github.com/uptrace/bun"
	"time"
)

type ActivationTokenRepository struct {
	db *bun.DB
}

func NewActivationTokenRepository(db *bun.DB) *ActivationTokenRepository {
	return &ActivationTokenRepository{db: db}
}

// ActivationToken
func (tr *ActivationTokenRepository) CreateAToken(ctx context.Context, aToken *tokens.ActivationToken,
	expiresAt time.Time) error {
	_, err := tr.db.NewInsert().
		Model(aToken).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (tr *ActivationTokenRepository) ActivationToken(ctx context.Context,
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

func (tr *ActivationTokenRepository) DeleteActivationToken(ctx context.Context,
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

func (tr *ActivationTokenRepository) SetUsedAt(ctx context.Context, userId int64) error {
	usedAt := time.Now()
	_, err := tr.db.NewUpdate().
		Model((*tokens.ActivationToken)(nil)).
		Set("used_at = ?", &usedAt).
		Where("user_id = ?", userId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
