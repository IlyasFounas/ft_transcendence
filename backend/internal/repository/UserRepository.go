package repository

import (
	// "fmt"
	"errors"
	"context"
	"github.com/uptrace/bun"
    "backend/internal/models"
	"backend/internal/tools"
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

func CreateUser(ctx context.Context, db *bun.DB, user *model.User) (*model.User, error) {
	_, err := db.NewInsert().
				Model(user).
				Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func LoginUser(ctx context.Context, db *bun.DB,
	login string, password string) (bool, int64, error) {

	if len(login) == 0 || len(password) == 0{
		return false, 0, nil
	}
	user := new(model.User)
	err := db.NewSelect().
		Model(user).
		Where("login = ?", login).
		Scan(ctx)
	if err != nil {
		return false, 0, err
	}
	if !tool.CheckPasswordHash(password, user.Password) {

	} else if !user.Is_verified {
		return false, user.ID, errors.New("The user has not verified his email")
	}
	return true, user.ID, nil
}

func GetUser(ctx context.Context, db *bun.DB, userId int64) (*model.User, error) {
	user := new(model.User)
	err := db.NewSelect().
			Model(user).
			Column("id", "name", "email", "is_verified","avatar", "created_at").
			Where("id = ?", userId).
			Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetName(ctx context.Context, db *bun.DB, userId int64) (*string, error) {
	user := new(model.User)
	err := db.NewSelect().
			Model(user).
			Column("name").
			Where("id = ?", userId).
			Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user.Username, nil
}

func GetLogin(ctx context.Context, db *bun.DB, userId int64) (*string, error) {
	user := new(model.User)
	err := db.NewSelect().
			Model(user).
			Column("login").
			Where("id = ?", userId).
			Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user.Login, nil
}

func GetEmail(ctx context.Context, db *bun.DB, userId int64) (*string, error) {
	user := new(model.User)
	err := db.NewSelect().
			Model(user).
			Column("email").
			Where("id = ?", userId).
			Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user.Email, nil
}

func GetOAuth(ctx context.Context, db *bun.DB, userId int64) (*oauth2.Token, error) {
	user := new(model.User)
	err := db.NewSelect().
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
func SetName(ctx context.Context, db *bun.DB, userId int64, name string) (bool, error) {
	if len(name) == 0 {
		return false, nil
	}
	user := &model.User{Username: name}
	_, err := db.NewUpdate().
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
func SetLogin(ctx context.Context, db *bun.DB, userId int64, login string) (bool, error) {
	if len(login) == 0 {
		return false, nil
	}
	user := &model.User{Login: login}
	_, err := db.NewUpdate().
		Model(user).
		Column("login").
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func SetAvatar(ctx context.Context, db *bun.DB, userId int64, avatar string) (bool, error) {
	if len(avatar) == 0 {
		return false, nil
	}
	user := &model.User{Avatar: &avatar}
	_, err := db.NewUpdate().
		Model(user).
		Column("avatar").
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func SetIsVerified(ctx context.Context, db *bun.DB, userId int64) error {
	user := &model.User{Is_verified: true}
	_, err := db.NewUpdate().
		Model(user).
		Column("is_verified").
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}