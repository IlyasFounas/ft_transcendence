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
        user, err := repository.CreateUser(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "state": "failed", "message": "database exection failed"})
            return 
        } else {
            aToken, err := repository.CreateAToken(c.Request.Context(), user, db)
            if err != nil {
                fmt.Printf("Token creation error message : %s\n", err)
                c.JSON(http.StatusInternalServerError, gin.H{
                    "state": "failed", "message": "token creation went wrong"})
                return
            }
            if err := tool.SendActivationEmail(c.PostForm("email"), *aToken); err != nil {
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
        state, userId, err := repository.LoginUser(c, db)
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
                    "message": "token creation failed, please retry"})
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
        userData, err := repository.GetUser(c, db)
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
        name, err := repository.GetName(c, db)
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
        login, err := repository.GetLogin(c, db)
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
        email, err := repository.GetEmail(c, db)
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
        oauth, err := repository.GetOAuth(c, db)
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
        state, err := repository.SetName(c, db)
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
        state, err := repository.SetLogin(c, db)
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
        state, err := repository.SetAvatar(c, db)
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
