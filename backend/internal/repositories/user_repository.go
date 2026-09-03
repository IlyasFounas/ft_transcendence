package repositories

import (
	"database/sql"
	"backend/internal/models/users"
	"backend/internal/tools/encryption"
	"context"
	"errors"
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
		return false, user.ID, errors.New("The user has not set a password")
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
		Column("id", "name", "email", "is_verified", "avatar", "created_at").
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
		Column("id", "name", "email", "is_verified", "avatar", "created_at").
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
	Column("id", "name", "email", "is_verified", "avatar", "created_at").
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
	Column("id", "name", "email", "is_verified", "avatar", "created_at").
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
// todo: set a limit to change the name (maybe 1 change per month)
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

// todo: set a limit to change the login (maybe 1 change per month)
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
