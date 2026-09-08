package repositories

import (
	"backend/internal/models/friends"
	"context"
	"github.com/uptrace/bun"
)

type FriendRepository struct {
	db *bun.DB
}

func NewFriendRepository(db *bun.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

func (r *FriendRepository) AddFriend(ctx context.Context, friend *friends.Friend) (*friends.Friend, error) {
	_, err := r.db.NewInsert().
		Model(friend).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return friend, nil
}