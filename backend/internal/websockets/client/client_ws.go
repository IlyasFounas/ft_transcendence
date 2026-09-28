package client_ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Connexion struct {
    Conn        *websocket.Conn
    Send        chan []byte
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
    deadline := time.Now().Add(5*time.Second)
    err := c.Conn.WriteControl(websocket.CloseMessage,
        websocket.FormatCloseMessage(websocket.CloseGoingAway, "server shutting down"),
        deadline)
    if err != nil {
        log.Printf("error : %s\n", err)
        return
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
