package user_repo

import (
	// "fmt"
	// "errors"
	"github.com/uptrace/bun"
    "github.com/gin-gonic/gin"
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

func CreateUser(c *gin.Context, db *bun.DB) (bool, error) {
	hashpassword, hashErr := tool.HashPassword(c.PostForm("password"))
	if hashErr != nil {
		return false, hashErr
	}
	user := &user.User{Login: c.PostForm("login"), 
		Password: hashpassword, 
		Username: c.PostForm("username"),
		Email: c.PostForm("email"),
		Is_verified: false}
	
	_, err := db.NewInsert().
				Model(user).
				Exec(c)
	if err != nil {
		return false, err
	}
	return true, nil
}

func LoginUser(c *gin.Context, db *bun.DB) (bool, int64, error) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	if len(login) == 0 || len(password) == 0{
		return false, 0, nil
	}
	user := new(user.User)
	err := db.NewSelect().
		Model(user).
		Where("login = ?", login).
		Scan(c)
	if err != nil {
		return false, 0, err
	}
	if !tool.CheckPasswordHash(password, user.Password) {
		return false, 0, nil
	}
	return true, user.ID, nil
}

func GetUser(c *gin.Context, db *bun.DB) (*user.User, error) {
	userId, cErr := tool.ConvertionUserId(c)
	if cErr != nil{
		return nil, cErr
	}
	user := new(user.User)
	err := db.NewSelect().
			Model(user).
			Column("id", "name", "email", "is_verified","avatar", "created_at").
			Where("id = ?", userId).
			Scan(c)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetName(c *gin.Context, db *bun.DB) (*string, error) {
	userId, cErr := tool.ConvertionUserId(c)
	if cErr != nil{
		return nil, cErr
	}
	user := new(user.User)
	err := db.NewSelect().
			Model(user).
			Column("name").
			Where("id = ?", userId).
			Scan(c)
	if err != nil {
		return nil, err
	}
	return &user.Username, nil
}

func GetLogin(c *gin.Context, db *bun.DB) (*string, error) {
	userId, cErr := tool.ConvertionUserId(c)
	if cErr != nil{
		return nil, cErr
	}
	user := new(user.User)
	err := db.NewSelect().
			Model(user).
			Column("login").
			Where("id = ?", userId).
			Scan(c)
	if err != nil {
		return nil, err
	}
	return &user.Login, nil
}

func GetEmail(c *gin.Context, db *bun.DB) (*string, error) {
	userId, cErr := tool.ConvertionUserId(c)
	if cErr != nil{
		return nil, cErr
	}
	user := new(user.User)
	err := db.NewSelect().
			Model(user).
			Column("email").
			Where("id = ?", userId).
			Scan(c)
	if err != nil {
		return nil, err
	}
	return &user.Email, nil
}

func GetOAuth(c *gin.Context, db *bun.DB) (*oauth2.Token, error) {
	userId, cErr := tool.ConvertionUserId(c)
	if cErr != nil{
		return nil, cErr
	}
	user := new(user.User)
	err := db.NewSelect().
			Model(user).
			Column("oauth").
			Where("id = ?", userId).
			Scan(c)
	if err != nil {
		return nil, err
	}
	return user.OAuth, nil
}

// todo: set a limit to change the name (maybe 1 change per month)
func SetName(c *gin.Context, db *bun.DB) (bool, error) {
	userId, cErr := tool.ConvertionUserId(c)
	if cErr != nil{
		return false, cErr
	}
	name := c.PostForm("name")

	if len(name) == 0 {
		return false, nil
	}
	user := &user.User{Username: name}
	_, err := db.NewUpdate().
		Model(user).
		Column("name").
		Where("id = ?", userId).
		Exec(c)
	if err != nil {
		return false, err
	}
	return true, nil
}

// todo: set a limit to change the login (maybe 1 change per month)
func SetLogin(c *gin.Context, db *bun.DB) (bool, error) {
	userId, cErr := tool.ConvertionUserId(c)
	if cErr != nil{
		return false, cErr
	}
	login := c.PostForm("login")

	if len(login) == 0 {
		return false, nil
	}
	user := &user.User{Login: login}
	_, err := db.NewUpdate().
		Model(user).
		Column("login").
		Where("id = ?", userId).
		Exec(c)
	if err != nil {
		return false, err
	}
	return true, nil
}

func SetAvatar(c *gin.Context, db *bun.DB) (bool, error) {
	userId, cErr := tool.ConvertionUserId(c)
	if cErr != nil{
		return false, cErr
	}
	avatar := c.PostForm("avatar")

	if len(avatar) == 0 {
		return false, nil
	}
	user := &user.User{Avatar: &avatar}
	_, err := db.NewUpdate().
		Model(user).
		Column("avatar").
		Where("id = ?", userId).
		Exec(c)
	if err != nil {
		return false, err
	}
	return true, nil
}