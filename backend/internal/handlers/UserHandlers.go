package handler

import (
    // EXTERNAL PACKAGE
    "os"
    "fmt"
    "time"
    "net/http"
	"backend/internal/tools"
	"github.com/uptrace/bun"
    "backend/internal/models"
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

func createActivationToken(c *gin.Context, db *bun.DB, expiresAt time.Time,
        userId int64) error {          
    tokenObj, tokenErr := tool.GenerateToken(userId, expiresAt, 
        []byte(os.Getenv("ACTIVATION_JWT_SECRET")))
    if tokenErr != nil {
        c.JSON(http.StatusUnauthorized,
            gin.H{"message": "Send a new activation link"})
        return tokenErr
    }
    err := repository.DeleteActivationToken(c.Request.Context(), db, userId) 
    if err != nil {
        c.JSON(http.StatusInternalServerError,
            gin.H{"message": "Send a new activation link"})
        return err
    }
    err = repository.CreateAToken(c.Request.Context(), db, tokenObj, expiresAt)
    if err != nil {
        fmt.Printf("Token creation error message : %s\n", err)
        c.JSON(http.StatusInternalServerError,
            gin.H{"message": "Send a new activation link"})
        return err
    }
    // we maybe want to ue goroutine for this especially for sending emails
    if err := tool.SendActivationEmail(c.PostForm("email"), tokenObj.Token); err != nil {
        c.JSON(http.StatusInternalServerError,
            gin.H{"message": "Send a new activation link"})
        return err
    }
    return nil
}

func CreateUser(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        hashpassword, hashErr := tool.HashPassword(c.PostForm("password"))
        if hashErr != nil {
            c.JSON(http.StatusInternalServerError,
                gin.H{"message": "Something went wrong"})
            return
        }
        userObj := &model.User{Login: c.PostForm("login"), 
            Password: hashpassword, 
            Username: c.PostForm("username"),
            Email: c.PostForm("email"),
            Is_verified: false}
        
        user, err := repository.CreateUser(c.Request.Context(), db, userObj)
        if err != nil {
            c.JSON(http.StatusInternalServerError,
                gin.H{"message": "The account creation failed"})
            return
        } else {
            expiresAt := time.Now().Add(24 * time.Hour)
            if err := createActivationToken(c, db, expiresAt, user.ID); err != nil {
                return 
            }
            c.JSON(http.StatusCreated,
                gin.H{"message": "The activation link has been sent on your email"})
        }
    }
}

func LoginUser(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        state, userId, err := repository.LoginUser(c.Request.Context(), db,
            c.PostForm("login"), c.PostForm("password"))
        if err != nil && userId != 0 {
            expiresAt := time.Now().Add(24 * time.Hour)
            if err := createActivationToken(c, db, expiresAt, userId); err != nil {
                return 
            }
            c.JSON(http.StatusInternalServerError,
                gin.H{"message": "The activation link has been sent on your email"})
            return
        } else if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
            return 
        }
        if !state {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "Try to login again"})
        } else {
            // generate the jwt token here
            expiresAt := time.Now().Add(15 * time.Minute)
            tokenObj, tokenErr := tool.GenerateToken(userId, expiresAt, 
                []byte(os.Getenv("JWT_SECRET")))
            if tokenErr != nil {
                c.JSON(http.StatusInternalServerError,
                    gin.H{"message": "token creation went wrong"})
                return
            }
            encodedKey, rTokenErr := CreateRefreshToken(c.Request.Context(), db, userId)
            if rTokenErr != nil {
                c.JSON(http.StatusInternalServerError,
                    gin.H{"message": "Try to login again"})
                return
            }
            c.JSON(http.StatusOK,
                gin.H{"Authorization": tokenObj.Token, "Refresh-Token": encodedKey})
        }
    }
}

// Getters
func GetUser(db *bun.DB) gin.HandlerFunc { //return almost all user infos
    return func(c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.JSON(http.StatusInternalServerError,
                gin.H{"message": "The userId is empty"})
            return
        }
        userData, err := repository.GetUser(c.Request.Context(), db, userId.(int64))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "The user doesn't exist"})
            return
        }
        if userData == nil {
            c.JSON(http.StatusNotFound,
                gin.H{"message": "The user doesn't exist"})
            return
        } else {
            c.JSON(http.StatusOK,
                gin.H{"message": "success", "user":  userData})
            return
        }
    }
}

func GetName(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(500,
                gin.H{"message": "The user id doesn't exist"})
        }
        name, err := repository.GetName(c.Request.Context(), db, userId.(int64))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
        }
        if name == nil {
            c.JSON(http.StatusNotFound, gin.H{"message": "failed"})
        } else {
            c.JSON(http.StatusOK,
                gin.H{"message": "success", "name":  name})
        }
    }
}

func GetLogin(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(http.StatusInternalServerError,
                gin.H{"message": "The user id doesn't exist"})
        }
        login, err := repository.GetLogin(c.Request.Context(), db, userId.(int64))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
        }
        if login == nil {
            c.JSON(http.StatusNotFound, gin.H{"message": "failed"})
        } else {
            c.JSON(http.StatusOK, gin.H{"message": "success", "login":  login})
        }
    }
}

func GetEmail(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(http.StatusInternalServerError,
                gin.H{"message": "The user id doesn't exist"})
        }
        email, err := repository.GetEmail(c.Request.Context(), db, userId.(int64))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
        }
        if email == nil {
            c.JSON(http.StatusNotFound, gin.H{"message": "the user have no email"})
        } else {
            c.JSON(http.StatusOK, gin.H{"message": "success", "email":  email})
        }
    }
}

func GetOAuth(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(http.StatusInternalServerError,
                gin.H{"message": "The user id doesn't exist"})
        }
        oauth, err := repository.GetOAuth(c.Request.Context(), db, userId.(int64))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
        }
        if oauth == nil {
            c.JSON(http.StatusNotFound, gin.H{"message":  "oauth not found"})
        } else {
            c.JSON(http.StatusOK, gin.H{
            "message": "success",
            "OAuth":  oauth})
        }
    }
}

// Setters
func SetName(db *bun.DB) gin.HandlerFunc {
    return func (c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(http.StatusInternalServerError,
                gin.H{"message": "The user id doesn't exist"})
        }
        state, err := repository.SetName(c.Request.Context(), db, userId.(int64), c.PostForm("name"))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{"message": "failed"})
        } else {
             c.JSON(http.StatusOK, gin.H{"message": "success"})
        }
    }
}

func SetLogin(db *bun.DB) gin.HandlerFunc {
    return func (c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(http.StatusInternalServerError,
                gin.H{"message": "The user id doesn't exist"})
        }
        state, err := repository.SetLogin(c.Request.Context(), db, userId.(int64), c.PostForm("login"))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "database function failed"})
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{"message": "failed"})
        } else {
             c.JSON(http.StatusOK, gin.H{"message": "success"})
        }
    }
}

func SetAvatar(db *bun.DB) gin.HandlerFunc {
    return func (c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(http.StatusInternalServerError,
                gin.H{"message": "The user id doesn't exist"})
        }
        state, err := repository.SetAvatar(c.Request.Context(), db, userId.(int64), c.PostForm("avatar"))
        if err != nil {
            c.JSON(http.StatusInternalServerError,
                gin.H{"message": "database function failed"})
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{"message": "failed"})
        } else {
             c.JSON(http.StatusOK, gin.H{"message": "success"})
        }
    }
}
