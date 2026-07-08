package user

import (
    // EXTERNAL PACKAGE
    "net/http"
	"github.com/uptrace/bun"
    "github.com/gin-gonic/gin"
    "backend/internal/repository"
)

// Public functions
func HomePage(c *gin.Context) {
    c.String(http.StatusOK, "Hello world from backend !")
}

func LoginUser(c *gin.Context) {
    // call the repostirories func later
    login := c.PostForm("login")
    password := c.PostForm("password")

    c.JSON(http.StatusOK, gin.H{"state": "succesfull connexion", "name": login, "password": password})
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
    // recup the JWT token later
    // call the repostirories func later, give it the id and check the token

    c.JSON(http.StatusOK, gin.H{"name": "ifounas", "password": "salut", "id_user": id_user})
}
