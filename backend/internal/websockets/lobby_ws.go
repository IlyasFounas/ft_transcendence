package ws

import (
	"gorilla/websocket"
)

type Upgrader struct {
	upgrader websocket.Upgrader
	userId int64
}

func WsLobby(c *gin.Context) {
    conn, err := Upgrader.upgrader.Upgrade(c.Writter(), c.Request(), nil)
    if err != nil {
        log.Println(err)
        return
    }
}