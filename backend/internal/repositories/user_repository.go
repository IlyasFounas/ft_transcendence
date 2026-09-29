package repositories

import (
	"database/sql"
	"backend/internal/models/users"
	"backend/internal/tools/encryption"
	"backend/internal/tools/rank_calculation"
	"context"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
)

/**
 * Repositories and migrations are the only places where
 * you directly interact with the database. Therefore, when working here,
 * you don't need to check if the database is working or not.
 *
 * For more details on how to perform CRUD operations (Create, Read, Update, Delete),
 * refer to the official documentation:
 * https://bun.uptrace.dev/guide/golang-orm.html#crud-operations
 */

// Leaderboard is a line of the leaderboard: the position of a player, his name and his xp
type Leaderboard struct {
	Position int64  `json:"position"`
	Username string `json:"username"`
	Rank string 	`json:"rank"`
	Total    int64  `json:"total"`
}

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *users.User) (*users.User, error) {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) LoginUser(ctx context.Context,
	login string, password string) (bool, int64, error) {
	if len(login) == 0 || len(password) == 0 {
		return false, 0, nil
	}
	user := new(users.User)
	err := r.db.NewSelect().
		Model(user).
		Where("login = ?", login).
		Scan(ctx)
	if err != nil {
		return false, 0, err
	}
	if user.Password == nil || !encryption.CheckPasswordHash(password, *user.Password) {
		return false, 0, errors.New("The user has not set a password")
	} else if !user.Is_verified {
		return false, user.ID, errors.New("The user has not verified his email")
	}
	return true, user.ID, nil
}

// Getters
func (r *UserRepository) User(ctx context.Context, userId int64) (*users.User, error) {
	user := new(users.User)
	err := r.db.NewSelect().
		Model(user).
		Column("id", "name", "email", "is_verified", "avatar", "xp", "rank", "created_at").
		Where("id = ?", userId).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Name(ctx context.Context, userId int64) (*string, error) {
	user := new(users.User)
	err := r.db.NewSelect().
		Model(user).
		Column("name").
		Where("id = ?", userId).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user.Username, nil
}

// 
func (r *UserRepository) ExistingName(ctx context.Context, name string) (bool, error) {
	exists, err := r.db.NewSelect().
		Model((*users.User)(nil)).
		Where("name = ?", name).
		Exists(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) Login(ctx context.Context, userId int64) (*string, error) {
	user := new(users.User)
	err := r.db.NewSelect().
		Model(user).
		Column("login").
		Where("id = ?", userId).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user.Login, nil
}

func (r *UserRepository) ExistingLogin(ctx context.Context, login string) (*string, error) {
	exists, err := r.db.NewSelect().
		Model((*users.User)(nil)).
		Where("login = ?", login).
		Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	return &login, nil
}

func (r *UserRepository) Email(ctx context.Context, userId int64) (*string, error) {
	user := new(users.User)
	err := r.db.NewSelect().
		Model(user).
		Column("email").
		Where("id = ?", userId).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user.Email, nil
}

func (r *UserRepository) UserByName(ctx context.Context,
	username string) (*users.User, error) {
	user := new(users.User)
	err := r.db.NewSelect().
		Model(user).
		Column("id", "name", "email", "is_verified", "avatar", "xp", "rank", "created_at").
		Where("name = ?", username).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UserByEmail(ctx context.Context,
	email string) (*users.User, error) {
	user := new(users.User)
	err := r.db.NewSelect().
	Model(user).
	Column("id", "name", "email", "is_verified", "avatar", "xp", "rank", "created_at").
	Where("email = ?", email).
	Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UserByGoogleId(ctx context.Context,
	googleId *string) (*users.User, error) {
	user := new(users.User)
	err := r.db.NewSelect().
	Model(user).
	Column("id", "name", "email", "is_verified", "avatar", "xp", "rank", "created_at").
	Where("google_id = ?", googleId).
	Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Setters
func (r *UserRepository) SetName(ctx context.Context, userId int64, name string) (bool, error) {
	if len(name) == 0 {
		return false, nil
	}
	user := &users.User{Username: name}
	_, err := r.db.NewUpdate().
		Model(user).
		Column("name").
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepository) SetGoogleId(ctx context.Context, userId int64, googleId *string) (bool, error) {
	if googleId == nil || *googleId == "" {
		return false, nil
	}
	user := &users.User{GoogleId: googleId}
	_, err := r.db.NewUpdate().
		Model(user).
		Column("google_id").
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepository) SetLogin(ctx context.Context, userId int64, login string) (bool, error) {
	if len(login) == 0 {
		return false, nil
	}
	user := &users.User{Login: &login}
	_, err := r.db.NewUpdate().
		Model(user).
		Column("login").
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepository) SetAvatar(ctx context.Context, userId int64, avatar string) (bool, error) {
	if len(avatar) == 0 {
		return false, nil
	}
	user := &users.User{Avatar: &avatar}
	_, err := r.db.NewUpdate().
		Model(user).
		Column("avatar").
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepository) SetIsVerified(ctx context.Context, userId int64) error {
	user := &users.User{Is_verified: true}
	_, err := r.db.NewUpdate().
		Model(user).
		Column("is_verified").
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Getter xp
func (r *UserRepository) XP(ctx context.Context, userId int64) (int64, error) {
	user := new(users.User)
	err := r.db.NewSelect().
		Model(user).
		Column("xp").
		Where("id = ?", userId).
		Scan(ctx)
	if err != nil {
		return 0, err
	}
	return user.XP, nil
}

// Setter xp: replaces the xp of a user (a negative xp is refused)
func (r *UserRepository) SetXP(ctx context.Context, userId int64, xp int64) (bool, error) {
	if xp < 0 {
		return false, nil
	}
	user := &users.User{XP: xp}
	_, err := r.db.NewUpdate().
		Model(user).
		Column("xp").
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Returns the id of every user (to update every player when the server starts, see Repositories)
func (r *UserRepository) AllIds(ctx context.Context) ([]int64, error) {
	var ids []int64
	err := r.db.NewSelect().
		Model((*users.User)(nil)).
		Column("id").
		Scan(ctx, &ids)
	return ids, err
}

// Update the xp of each user using the sum of the base_score of all his guesses in all his FINISHED games 
func (r *UserRepository) UpdateXP(ctx context.Context, userIds []int64) error {
	if len(userIds) == 0 {
		return nil
	}
	var totals []struct {
		UserID int64 `bun:"user_id"`
		XP     int64 `bun:"xp"`
	}
	err := finishedGamesGuesses(r.db).
		ColumnExpr("g.user_id").
		ColumnExpr("SUM(g.base_score) AS xp").
		Where("g.user_id IN (?)", bun.In(userIds)).
		Group("g.user_id").
		Scan(ctx, &totals)
	if err != nil {
		return err
	}
	xp := make(map[int64]int64, len(totals))
	for _, total := range totals {
		xp[total.UserID] = total.XP
	}
	var errs []error
	for _, userId := range userIds {
		// a user without any guess in a finished game gets 0
		if _, err := r.SetXP(ctx, userId, xp[userId]); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Update the rank of each user from his current xp and, every user is tried even if one fails.
func (r *UserRepository) UpdateRank(ctx context.Context, userIds []int64) error {
	var errs []error
	for _, userId := range userIds {
		xp, err := r.XP(ctx, userId)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if _, err := r.SetRank(ctx, userId, rank_calculation.CalculateRank(xp)); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Leaderboard returns the 100 first players ranked by xp
func (r *UserRepository) Leaderboard(ctx context.Context) ([]Leaderboard, error) {
	var rows []Leaderboard
	err := r.db.NewSelect().
		Model((*users.User)(nil)).
		ColumnExpr("RANK() OVER (ORDER BY u.xp DESC) AS position").
		ColumnExpr("u.name AS username").
		ColumnExpr("u.rank AS rank").
		ColumnExpr("u.xp AS total").
		OrderExpr("u.xp DESC, u.name ASC").
		Limit(100).
		Scan(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	return rows, nil
}

// UserLeaderboard returns the position, name and xp of the user userId.
// The rank is computed over ALL the users first, then filtered (filtering first would always rank him 1st)
func (r *UserRepository) UserLeaderboard(ctx context.Context, userId int64) ([]Leaderboard, error) {
	ranked := r.db.NewSelect().
		Model((*users.User)(nil)).
		ColumnExpr("u.id").
		ColumnExpr("RANK() OVER (ORDER BY u.xp DESC) AS position").
		ColumnExpr("u.name AS username").
		ColumnExpr("u.rank AS rank").
		ColumnExpr("u.xp AS total")
	var rows []Leaderboard
	err := r.db.NewSelect().
		TableExpr("(?) AS ranked", ranked).
		ColumnExpr("ranked.position, ranked.username, ranked.rank, ranked.total").
		Where("ranked.id = ?", userId).
		Scan(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard info for user %d: %w", userId, err)
	}
	return rows, nil
}

func (r *UserRepository) SetRank(ctx context.Context, userId int64, rank string) (bool, error) {
	user := &users.User{Rank: rank}
	_, err := r.db.NewUpdate().
		Model(user).
		Column("rank").
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}