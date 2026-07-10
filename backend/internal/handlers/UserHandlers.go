package user

import (
    // EXTERNAL PACKAGE
    "net/http"
	"github.com/uptrace/bun"
    "github.com/gin-gonic/gin"
    "backend/internal/repository"
)

/**
 * README:
 * Handlers make the bridge between the router and the repository functions.
 * It's important to split different features into different files to improve visibility
 * and maintainability.
 *
 * Basically, the handlers call the repository functions, manage any errors,
 * and return the result to the frontend (in JSON, plain text, etc.).
 */

// Auth functions
func LoginUser(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        state, id_user, err := user_repo.LoginUser(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{
                "state": "login failed","id_user": id_user})
        } else {
            c.JSON(http.StatusOK, gin.H{
                "state": "Logged in successfully","id_user": id_user})
        }
    }
}

func CreateUser(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        state, err := user_repo.CreateUser(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "state": "database function failed"})
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{
                "state": "creation failed"})
        } else {
            c.JSON(http.StatusCreated, gin.H{
                "state": "successful creation"})
        }
    }
}

// Getters
func GetUser(db *bun.DB) gin.HandlerFunc { //return almost all user infos
    return func(c *gin.Context) {
        userData, err := user_repo.GetUser(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
        }
        if userData == nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        } else {
            c.JSON(http.StatusOK, gin.H{
            "state": "success",
            "user":  userData,
        })
        }
    }
}

func GetName(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        name, err := user_repo.GetName(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
        }
        if name == nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "name not found"})
        } else {
            c.JSON(http.StatusOK, gin.H{
            "name":  name})
        }
    }
}

func GetLogin(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        login, err := user_repo.GetLogin(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
        }
        if login == nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "login not found"})
        } else {
            c.JSON(http.StatusOK, gin.H{
            "login":  login})
        }
    }
}

func GetEmail(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        email, err := user_repo.GetEmail(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
        }
        if email == nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
        } else {
            c.JSON(http.StatusOK, gin.H{
            "email":  email})
        }
    }
}

func GetOAuth(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        oauth, err := user_repo.GetOAuth(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
        }
        if oauth == nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "oauth not found"})
        } else {
            c.JSON(http.StatusOK, gin.H{
            "OAuth":  oauth})
        }
    }
}

// Setters
func SetName(db *bun.DB) gin.HandlerFunc {
    return func (c *gin.Context) {
        state, err := user_repo.SetName(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{
                "state": "modification didn't occur"})
        } else {
             c.JSON(http.StatusOK, gin.H{
                "state": "successful modification"})
        }
    }
}

func SetLogin(db *bun.DB) gin.HandlerFunc {
    return func (c *gin.Context) {
        state, err := user_repo.SetLogin(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{
                "state": "modification didn't occur"})
        } else {
             c.JSON(http.StatusOK, gin.H{
                "state": "successful modification"})
        }
    }
}

func SetAvatar(db *bun.DB) gin.HandlerFunc {
    return func (c *gin.Context) {
        state, err := user_repo.SetAvatar(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{
                "state": "modification didn't occur"})
        } else {
             c.JSON(http.StatusOK, gin.H{
                "state": "successful modification"})
        }
    }
}
