package handler

import (
    // EXTERNAL PACKAGE
    "os"
    "fmt"
    "time"
    "strconv"
    "net/http"
	"backend/internal/tools"
	"github.com/uptrace/bun"
    "backend/internal/models"
    "github.com/gin-gonic/gin"
    "backend/internal/repository"
    "github.com/golang-jwt/jwt/v5"
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

type UserClaims struct {
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// Auth functions
func CreateUser(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        hashpassword, hashErr := tool.HashPassword(c.PostForm("password"))
        if hashErr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "state": "failed", "message": "password encryption failed"})
            return
        }
        userObj := &model.User{Login: c.PostForm("login"), 
            Password: hashpassword, 
            Username: c.PostForm("username"),
            Email: c.PostForm("email"),
            Is_verified: false}
        
        user, err := repository.CreateUser(c.Request.Context(), db, userObj)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "state": "failed", "message": "database exection failed"})
            return
        } else {
            expiresAt := time.Now().Add(24 * time.Hour)
            tokenObj, tokenErr := tool.GenerateToken(user.ID, expiresAt)
            if tokenErr != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
                    "state": "failed", "message": "token creation went wrong"})
                return
            }
            err := repository.CreateAToken(c.Request.Context(), db, tokenObj, expiresAt)
            if err != nil {
                fmt.Printf("Token creation error message : %s\n", err)
                c.JSON(http.StatusInternalServerError, gin.H{
                    "state": "failed", "message": "token creation went wrong"})
                return
            }
            if err := tool.SendActivationEmail(c.PostForm("email"), tokenObj.Token); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
                    "state": "failed", "message": "link activation sending went wrong"})
                return
            }
            c.JSON(http.StatusCreated, gin.H{
                "state": "success", "message": "you just receive an activation link on your email"})
        }
    }
}

func LoginUser(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        state, userId, err := repository.LoginUser(c.Request.Context(), db,
            c.PostForm("login"), c.PostForm("password"))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
            return 
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{
                "state": "failed",
                "message": "The login or the password is wrong"})
        } else {
            // generate the jwt token here
            claims := UserClaims{
                Role: "user",
                RegisteredClaims: jwt.RegisteredClaims{
                    Subject: strconv.FormatInt(userId, 10),
                    ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
                    IssuedAt:  jwt.NewNumericDate(time.Now()),
                    Issuer:    "golang-api",
                },
            }
            token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
            tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
                    "state": "failed",
                    "message": "token creation failed, please retry to login"})
                return 
            }
            c.JSON(http.StatusOK, gin.H{
                "state": "success",
                "Authorization": tokenString})
        }
    }
}

// Getters
func GetUser(db *bun.DB) gin.HandlerFunc { //return almost all user infos
    return func(c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(500, gin.H{"state": "failed",
                "error_message": "The user id doesn't exist"})
        }
        userData, err := repository.GetUser(c.Request.Context(), db, userId.(int64))
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"state": "failed", "error_message": err})
        }
        if userData == nil {
            c.JSON(http.StatusNotFound, gin.H{"state": "failed", "message": "The user doesn't exist"})
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
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(500, gin.H{"state": "failed",
                "error_message": "The user id doesn't exist"})
        }
        name, err := repository.GetName(c.Request.Context(), db, userId.(int64))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
        }
        if name == nil {
            c.JSON(http.StatusNotFound, gin.H{"state": "failed"})
        } else {
            c.JSON(http.StatusOK, gin.H{
            "state": "success",
            "name":  name})
        }
    }
}

func GetLogin(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(500, gin.H{"state": "failed",
                "error_message": "The user id doesn't exist"})
        }
        login, err := repository.GetLogin(c.Request.Context(), db, userId.(int64))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
        }
        if login == nil {
            c.JSON(http.StatusNotFound, gin.H{"stsate": "failed"})
        } else {
            c.JSON(http.StatusOK, gin.H{
            "state": "success",
            "login":  login})
        }
    }
}

func GetEmail(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(500, gin.H{"state": "failed",
                "error_message": "The user id doesn't exist"})
        }
        email, err := repository.GetEmail(c.Request.Context(), db, userId.(int64))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
        }
        if email == nil {
            c.JSON(http.StatusNotFound, gin.H{"state": "failed", "message": "the user have no email"})
        } else {
            c.JSON(http.StatusOK, gin.H{
            "state": "success",
            "email":  email})
        }
    }
}

func GetOAuth(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(500, gin.H{"state": "failed",
                "error_message": "The user id doesn't exist"})
        }
        oauth, err := repository.GetOAuth(c.Request.Context(), db, userId.(int64))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
        }
        if oauth == nil {
            c.JSON(http.StatusNotFound, gin.H{"state": "failed", "message":  "oauth not found"})
        } else {
            c.JSON(http.StatusOK, gin.H{
            "state": "success",
            "OAuth":  oauth})
        }
    }
}

// Setters
func SetName(db *bun.DB) gin.HandlerFunc {
    return func (c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(500, gin.H{"state": "failed",
                "error_message": "The user id doesn't exist"})
        }
        state, err := repository.SetName(c.Request.Context(), db, userId.(int64), c.PostForm("name"))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "failed"})
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{
                "state": "failed"})
        } else {
             c.JSON(http.StatusOK, gin.H{
                "state": "success"})
        }
    }
}

func SetLogin(db *bun.DB) gin.HandlerFunc {
    return func (c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(500, gin.H{"state": "failed",
                "error_message": "The user id doesn't exist"})
        }
        state, err := repository.SetLogin(c.Request.Context(), db, userId.(int64), c.PostForm("login"))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{
                "state": "failed"})
        } else {
             c.JSON(http.StatusOK, gin.H{
                "state": "success"})
        }
    }
}

func SetAvatar(db *bun.DB) gin.HandlerFunc {
    return func (c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.AbortWithStatusJSON(500, gin.H{"state": "failed",
                "error_message": "The user id doesn't exist"})
        }
        state, err := repository.SetAvatar(c.Request.Context(), db, userId.(int64), c.PostForm("avatar"))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{
                "state": "failed"})
        } else {
             c.JSON(http.StatusOK, gin.H{
                "state": "success"})
        }
    }
}
