package lobby_ws

import (
    "log"
    "net/http"
    "strconv"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"backend/internal/websockets/hub"
)

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {return true},
}

type ClientWs struct {
    conn    *websocket.Conn
    send    chan []byte
    userId  int64
    lobbyId int64
    gameId  int64
    hub  *hub_ws.HubWs
}

func NewClient(hub *hub_ws.HubWs, conn *websocket.Conn,
    userId int64, lobbyId int64, gameId int64) *ClientWs {
    return &ClientWs{conn: conn,
    send: make(chan []byte, 256),
    userId: userId,
    lobbyId: lobbyId,
    gameId: gameId,
    hub: hub}
}

func (c *ClientWs) WritePump() {
    defer c.conn.Close()
    for msg := range c.send {
        if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            return
        }
    }
}

func (c *ClientWs) ReadPump() {
    defer func() { c.hub.Unregister <- c }()
    defer c.conn.Close()
    for {
        msgType, data, err := c.conn.ReadMessage()
        if err != nil {
            break // the connection ended up
        }
        if msgType == websocket.TextMessage {
            // TODO : give data to the HUB
            c.send <- data
        }
    }
}

func WsJoinLobby(hub *hub_ws.HubWs) gin.HandlerFunc {
	return func(c *gin.Context) {
        // userId, ok := c.Get("userId")
        // if !ok {
        //     c.JSON(http.StatusUnauthorized,
        //         gin.H{"message": "The userId is empty"})
        //     return
        // }
        userId, err := strconv.ParseInt(c.Query("userId"), 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "invalid userId"})
            return
        }
        lobbyId, err := strconv.ParseInt(c.Query("lobbyId"), 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "invalid lobbyId"})
            return
        }
        gameId, err := strconv.ParseInt(c.Query("gameId"), 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "invalid gameId"})
            return
        }

        // Check if the player is in the right lobby
        // We pass to an HTTP request to a websocket connexion
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            log.Println(err)
            return
        }
        // defer conn.Close()
        client := NewClient(hub, conn, userId, lobbyId, gameId)
        hub.Register <- client
        go client.WritePump()
        client.ReadPump()
    }
}