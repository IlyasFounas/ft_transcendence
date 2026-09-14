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

func (r *FriendRepository) RejectRequest(ctx context.Context, userId int64, rejectID int64) (bool, error) {
	result, err := r.db.NewDelete().
		Model((*friends.Friend)(nil)).
		Where("(requester_id = ? AND receiver_id = ?) AND status = ?", rejectID, userId, "pending").
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

func (r *FriendRepository) CancelRequest(ctx context.Context, userId int64, cancelID int64) (bool, error) {
	result, err := r.db.NewDelete().
		Model((*friends.Friend)(nil)).
		Where("(requester_id = ? AND receiver_id = ?) AND status = ?", userId, cancelID, "pending").
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

func (r *FriendRepository) RemoveFriend(ctx context.Context, userId int64, removeID int64) (bool, error) {
	result, err := r.db.NewDelete().
		Model((*friends.Friend)(nil)).
		Where("(requester_id = ? AND receiver_id = ?) AND status = ?", userId, removeID, "accepted").
		WhereOr("(requester_id = ? AND receiver_id = ?) AND status = ?", removeID, userId, "accepted").
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

func (r *FriendRepository) Friends(ctx context.Context, userId int64) ([]*friends.FriendResponse, error) {
    var friendList []*friends.FriendResponse
    err := r.db.NewSelect().
        TableExpr("friends AS f").
        Join(`
            JOIN users AS u ON u.id =
            CASE
                WHEN f.requester_id = ? THEN f.receiver_id
                ELSE f.requester_id
            END
        `, userId).
        ColumnExpr(`
            CASE
                WHEN f.requester_id = ? THEN f.receiver_id
                ELSE f.requester_id
            END AS user_id
        `, userId).
        Column("f.id", "u.name", "u.avatar").
		Where("(f.requester_id = ? OR f.receiver_id = ?)", userId, userId).
		Where("f.status = ?", "accepted").
        Order("f.created_at DESC").
        Scan(ctx, &friendList)

    if err != nil {
        return nil, err
    }
    return friendList, nil
}

func (r *FriendRepository) FriendRequests(ctx context.Context, userId int64) ([]*friends.FriendResponse, error) {
	var friendRequests []*friends.FriendResponse
	err := r.db.NewSelect().
		TableExpr("friends AS f").
        Join("JOIN users AS u ON u.id = f.requester_id").
		Column("f.id", "u.name", "u.avatar", "f.created_at").
		ColumnExpr("f.requester_id AS user_id").
		Where("f.receiver_id = ?", userId).
		Where("f.status = ?", "pending").
		Order("f.created_at DESC").
		Scan(ctx, &friendRequests)

	if err != nil {
		return nil, err
	}
	return friendRequests, nil
}

func (r *FriendRepository) FriendPending(ctx context.Context, userId int64) ([]*friends.FriendResponse, error) {
	var friendPending []*friends.FriendResponse
	err := r.db.NewSelect().
		TableExpr("friends AS f").
        Join("JOIN users AS u ON u.id = f.receiver_id").
		Column("f.id", "u.name", "u.avatar", "f.created_at").
		ColumnExpr("f.receiver_id AS user_id").
		Where("f.requester_id = ?", userId).
		Where("f.status = ?", "pending").
		Order("f.created_at DESC").
		Scan(ctx, &friendPending)

	if err != nil {
		return nil, err
	}
	return friendPending, nil
}