package user_repo

import (
	"fmt"
	"github.com/uptrace/bun"
    "github.com/gin-gonic/gin"
    "backend/internal/models"
	"backend/internal/tools"
)

func CreateUser(c *gin.Context, db *bun.DB) bool {
	hashpassword, hashErr := bcrypt.HashPassword(c.PostForm("password"))
	if hashErr != nil {
		return false
	}
	user := &user.User{Login: c.PostForm("login"), 
		Password: hashpassword, 
		Username: c.PostForm("username"),
		Email: c.PostForm("email"),
		Is_verified: false}
	_, err := db.NewInsert().Model(user).Exec(c)
	if err != nil {
		return false
	}
	return true
}

func LoginUser(c *gin.Context, db *bun.DB) (bool, int64) {
	login := c.PostForm("login")
	password := c.PostForm("password")
	user := new(user.User)
	err := db.NewSelect().
		Model(user).
		Where("login = ?", login).
		Scan(c)
	if err != nil {
		return false, 0
	}
	fmt.Printf("user login : %s\n", user.Login)
	if !bcrypt.CheckPasswordHash(password, user.Password) {
		return false, 0
	}
	return true, user.ID
}