package lobby_ws

import (
	"encoding/json"
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

type ClientWs struct {
    conn    *websocket.Conn
    send    chan []byte
    userId  int64
    lobbyId int64
    gameId  int64
    hub  *HubWs
}

func NewClient(hub *HubWs, conn *websocket.Conn,
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
    defer func() { c.hub.unregister <- c }()
    defer c.conn.Close()
    for {
        msgType, data, err := c.conn.ReadMessage()
        if err != nil {
            log.Printf("error : %s\n", err)
            break // the connection ended up
        }
        if msgType == websocket.TextMessage {
            // log.Printf("client data : %s\n", data)
            var serializedData DataInfos
            err := json.Unmarshal(data, &serializedData)
            if err != nil {
                log.Printf("error : %s\n", err)
                continue 
            }
            serializedData.Client = c
            c.hub.action <- &serializedData
        }
    }
}

func WsJoinLobby(hub *HubWs) gin.HandlerFunc {
	return func(c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "The userId is empty"})
            return
        }
        // Check if the player is in the right lobby
        lobby, err := hub.lobbyRepo.ActualLobby(c.Request.Context(), userId.(int64))
        if err != nil {
            log.Printf("error : %s\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"message": "could not verify lobby"})
            return
        }
        if lobby == nil {
            // Previously this returned with no response at all: 
            // the client's upgrade request was simply left hanging (neither a WS handshake nor an HTTP error)
            // so this message avoids a hung connection
            c.JSON(http.StatusForbidden, gin.H{"message": "you are not in a lobby"})
            return
        }
        // We pass to an HTTP request to a websocket connexion
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            log.Printf("error : %s\n", err)
            return
        }
        client := NewClient(hub, conn, userId.(int64), lobby.ID, lobby.GameId)
        client.hub.register <- client
        go client.WritePump()
        client.ReadPump()
    }
}