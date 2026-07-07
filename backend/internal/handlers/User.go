package user

import (
    // EXTERNAL PACKAGE
    "net/http"
    "github.com/gin-gonic/gin"
)

func UserRoute(r *gin.Engine) {
    r.GET("/user_test/:name", func(c *gin.Context) {
        name := c.Param("name")
        message := "welcome on your porfile " + name 
        c.String(http.StatusOK, message)
    })
}

func GameRoute(r *gin.Engine) {
    r.GET("/game", func(c *gin.Context) {
        player_name := c.DefaultQuery("player_name", "Guest mode")
        c.String(http.StatusOK, "The game will launch %s", player_name)
    })
}