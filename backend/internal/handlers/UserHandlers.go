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

// Public functions
func HomePage(c *gin.Context) {
    c.String(http.StatusOK, "Hello world from backend !")
}

func LoginUser(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        state, id_user := user_repo.LoginUser(c, db)
        if !state {
            c.JSON(http.StatusOK, gin.H{"state": "login failed", "id_user": id_user})
        } else {
            c.JSON(http.StatusOK, gin.H{"state": "Logged in successfully", "id_user": id_user})
        }
    }
}

func CreateUser(db *bun.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        state := user_repo.CreateUser(c, db)
        if !state {
            c.JSON(http.StatusOK, gin.H{"state": "creation failed"})
        } else {
            c.JSON(http.StatusOK, gin.H{"state": "successful creation"})
        }
    }
}

// Private functions
func GetUser(c *gin.Context) {
    id_user := c.Param("id")
    // recup the JWT token (JSON WEB TOKEN) later
    // call the repostirories func later, give it the id and check the token

    c.JSON(http.StatusOK, gin.H{"name": "ifounas", "password": "salut", "id_user": id_user})
}
