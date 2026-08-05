package repositories

import (
	"backend/internal/models/users"
	"backend/internal/tools"
	"context"
	"errors"
	"github.com/uptrace/bun"
	"golang.org/x/oauth2"
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
	if !tool.CheckPasswordHash(password, user.Password) {
		// TODO: add something like "return false, user.ID, nil" ??
	} else if !user.Is_verified {
		return false, user.ID, errors.New("The user has not verified his email")
	}
	return true, user.ID, nil
}

func (r *UserRepository) GetUser(ctx context.Context, userId int64) (*users.User, error) {
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

func (r *UserRepository) GetName(ctx context.Context, userId int64) (*string, error) {
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

func (r *UserRepository) GetLogin(ctx context.Context, userId int64) (*string, error) {
	user := new(users.User)
	err := r.db.NewSelect().
		Model(user).
		Column("login").
		Where("id = ?", userId).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user.Login, nil
}

func (r *UserRepository) GetEmail(ctx context.Context, userId int64) (*string, error) {
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

func (r *UserRepository) GetOAuth(ctx context.Context, userId int64) (*oauth2.Token, error) {
	user := new(users.User)
	err := r.db.NewSelect().
		Model(user).
		Column("oauth").
		Where("id = ?", userId).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user.OAuth, nil
}

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

// todo: set a limit to change the login (maybe 1 change per month)
func (r *UserRepository) SetLogin(ctx context.Context, userId int64, login string) (bool, error) {
	if len(login) == 0 {
		return false, nil
	}
	user := &users.User{Login: login}
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
