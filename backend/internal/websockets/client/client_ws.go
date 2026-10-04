package client_ws

import (
	"errors"
	"log"
	"net"
	"time"
    "sync"
	"github.com/gorilla/websocket"
)
const (
    writeWait = 10 * time.Second
    readWait = 60 * time.Second
    pingPeriod = 54 * time.Second
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

func (c *Connexion) WritePump(wg *sync.WaitGroup) {
	ticker := time.NewTicker(pingPeriod)
    defer wg.Done()
    defer c.Conn.Close()
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
                log.Printf("error : %s\n", err)
                return
            }
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
        case msg, ok := <- c.Send:
            if !ok {    
                deadline := time.Now().Add(2*time.Second)
                err := c.Conn.WriteControl(websocket.CloseMessage,
                    websocket.FormatCloseMessage(websocket.CloseGoingAway,
                        "the connexion is shutting down"), deadline)
                // the connection is usually already closed by ReadPump (the client left): not an error
                if err != nil && !errors.Is(err, net.ErrClosed) {
                    log.Printf("error : %s\n", err)
                }
                return
            }
            if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
                log.Printf("error : %s\n", err)
                return
            }
            if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
                log.Printf("error : %s\n", err)
                return
            }
        }
    }
}

func (c *Connexion) ReadPump(action func([]byte), unregister func(), wg *sync.WaitGroup) {
	defer wg.Done()
    defer unregister()
    defer c.Conn.Close()
    c.Conn.SetReadLimit(4096)
    if err := c.Conn.SetReadDeadline(time.Now().Add(readWait)); err != nil {
        log.Printf("error : %s\n", err)
        return
    }
    c.Conn.SetPongHandler(func(string) error {
        c.Conn.SetReadDeadline(time.Now().Add(readWait)); return nil })
    for {
        msgType, data, err := c.Conn.ReadMessage()
        if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway,
                websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
            break // the connection ended up
        }
        if msgType == websocket.TextMessage {
            action(data)
        }
    }
}