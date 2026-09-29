package client_ws

import (
	"log"
	"time"
    "sync"
	"github.com/gorilla/websocket"
)

const writeWait = 10 * time.Second

type Connexion struct {
    Conn        *websocket.Conn
    Send        chan []byte
}

func NewConnexion(conn *websocket.Conn) *Connexion {
    return &Connexion{Conn: conn,
    Send: make(chan []byte, 256),
    }
}

func (c *Connexion) WritePump(wg *sync.WaitGroup) {
	defer wg.Done()
    defer c.Conn.Close()
    for msg := range c.Send {
        if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
            log.Printf("error : %s\n", err)
            return
        }
        if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            log.Printf("error : %s\n", err)
            return
        }
    }
    deadline := time.Now().Add(2*time.Second)
    err := c.Conn.WriteControl(websocket.CloseMessage,
        websocket.FormatCloseMessage(websocket.CloseGoingAway, "the connexion is shutting down"),
        deadline)
    if err != nil {
        log.Printf("error : %s\n", err)
    }
}

func (c *Connexion) ReadPump(action func([]byte), unregister func(), wg *sync.WaitGroup) {
	defer wg.Done()
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
