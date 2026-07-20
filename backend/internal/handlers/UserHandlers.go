package user

import (
    // EXTERNAL PACKAGE
    "os"
    // "fmt"
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
        state, err := user_repo.CreateUser(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "state": "database function failed"})
            return 
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{
                "state": "creation failed"})
            return 
        } else {
            if err := tool.SendEmail(c.PostForm("email")); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
                    "state": "failed", "error_message": "code verification sending went wrong"})
                return 
            }
            c.JSON(http.StatusCreated, gin.H{
                "state": "successful creation"})
        }
    }
}

func LoginUser(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        state, userId, err := user_repo.LoginUser(c, db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"state": "database function failed"})
            return 
        }
        if !state {
            c.JSON(http.StatusOK, gin.H{
                "state": "login failed"})
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
                    "state": "token creation failed, please retry"})
                return 
            }
            c.JSON(http.StatusOK, gin.H{
                "state": "Logged in successfully",
                "Authorization": tokenString})
        }
    }
}

// Getters
func GetUser(db *bun.DB) gin.HandlerFunc { //return almost all user infos
    return func(c *gin.Context) {
        userData, err := user_repo.GetUser(c, db)
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"error": "Convertion failed", "error_message": err})
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
