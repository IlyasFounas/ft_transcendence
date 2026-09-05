package lobby_ws

import (
    "log"
    "net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {return true},
}

func WsLobby() gin.HandlerFunc {
	return func(c *gin.Context) {
        // userId, ok := c.Get("userId")
        // if !ok {
        //     c.JSON(http.StatusUnauthorized,
        //         gin.H{"message": "The userId is empty"})
        //     return
        // }
        // We pass to an HTTP request to a websocket connexion
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            log.Println(err)
            return
        }
        defer conn.Close()
        for {
            msgType, data, err := conn.ReadMessage()
            if err != nil {
                break // the connection ended up
            }
            if msgType == websocket.TextMessage {
                log.Printf("client request : %s\n", string(data))
            }
        }
    }
}