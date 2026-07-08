package user_repo

import (
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
		Email: c.PostForm("email")}
	_, err := db.NewInsert().Model(user).Exec(c)
	if err != nil {
		return false
	}
	return true
}

func LoginUser(c *gin.Context, db *bun.DB) bool {
	hashpassword, hashErr := bcrypt.HashPassword(c.PostForm("password"))
	if hashErr != nil {
		return false
	}
	user := &user.User{Login: c.PostForm("login"), 
		Password: hashpassword, 
		Username: c.PostForm("username"),
		Email: c.PostForm("email")}
	_, err := db.NewInsert().Model(user).Exec(c)
	if err != nil {
		return false
	}
	return true
}