package client_ws

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {return true},
}

type Connexion struct {
    Conn    *websocket.Conn
    Send    chan []byte
}

func NewConnexion(conn *websocket.Conn) *Connexion {
    return &Connexion{Conn: conn,
    Send: make(chan []byte, 256),
    }
}

func (c *Connexion) WritePump() {
    defer c.Conn.Close()
    for msg := range c.Send {
        if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            return
        }
    }
}

func (c *Connexion) ReadPump(action func([]byte), unregister func()) {
    defer unregister()
    defer c.Conn.Close()
    for {
        msgType, data, err := c.Conn.ReadMessage()
        if err != nil {
            log.Printf("error : %s\n", err)
            break // the connection ended up
        }
        if msgType == websocket.TextMessage {
            action(data)
        }
    }
}

// type ClientHub struct {
//     connx     *Connexion
//     userId  int64
//     lobbyId int64
//     gameId  int64
// }


// func NewClientHub(hub any, connx *Connexion,
//     userId int64, lobbyId int64, gameId int64) *ClientHub {
//     return &ClientHub{connx: connx,
//     userId: userId,
//     lobbyId: lobbyId,
//     gameId: gameId,
//     hub: hub}
// }