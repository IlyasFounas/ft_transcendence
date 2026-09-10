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

func (r *FriendRepository) AcceptRequest(ctx context.Context, userId int64, requesterID int64) (bool, error) {
	friend := new(friends.Friend)
	result, err := r.db.NewUpdate().
		Model(friend).
		Set("status = ?", "accepted").
		Where("receiver_id = ?", userId).
		Where("requester_id = ?", requesterID).
		Where("status = ?", "pending").
		Exec(ctx)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, nil
	}
	return true, nil
}

// func (r *FriendRepository) Friends(ctx context.Context, userId int64) ([]*friends.Friend, error) {
// 	var friends []*friends.Friend
// 	err := r.db.NewSelect().
// 		Model(&friends).
// 		Relation("Requester", func(q *bun.SelectQuery) *bun.SelectQuery {
//         	return q.Column("id", "name", "avatar")
// 		}).
// 		Relation("Receiver", func(q *bun.SelectQuery) *bun.SelectQuery {
//         	return q.Column("id", "name", "avatar")
// 		}).
// 		Where("(requester_id = ? OR receiver_id = ?) AND status = ?", userId, userId, "accepted").
// 		Order("created_at DESC").
// 		Scan(ctx)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return friends, nil
// }

func (r *FriendRepository) FriendRequests(ctx context.Context, userId int64) ([]*friends.Friend, error) {
	var friendRequests []*friends.Friend
	err := r.db.NewSelect().
		Model(&friendRequests).
		Column("id", "created_at").
		Relation("Requester", func(q *bun.SelectQuery) *bun.SelectQuery {
        	return q.Column("id", "name", "avatar")
		}).
		Where("receiver_id = ?", userId).
		Where("status = ?", "pending").
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return friendRequests, nil
}