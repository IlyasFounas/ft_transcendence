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

func (r *FriendRepository) FriendRequests(ctx context.Context, userId int64) ([]*friends.Friend, error) {
	var friendRequests []*friends.Friend
	err := r.db.NewSelect().
		Model(&friendRequests).
		Where("receiver_id = ?", userId).
		Where("status = ?", "pending").
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return friendRequests, nil
}