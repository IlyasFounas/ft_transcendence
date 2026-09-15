package chat_ws

import (
	// "log"
	// "context"
	// "strings"
	// "backend/internal/models/chats"
	"encoding/json"
	"backend/internal/repositories"
	"backend/internal/websockets/client"
	"backend/internal/websockets/handlers"
)

type Chat struct {
	clients map[int64]map[*ClientChat]bool
	unregister chan *ClientChat
	register chan *ClientChat
	action chan *DataInfos
	chatRepo *repositories.ChatRepository
	chatHandler *ws_handlers.ChatHandler
	userRepo *repositories.UserRepository
}

type ClientChat struct {
    connx     *client_ws.Connexion
    userId    int64
    targetId  int64
}


type DataInfos struct {
	Action	string
	Payload json.RawMessage
	Client	*ClientChat
}

func NewChat(
	repo *repositories.Repositories,
	wsHandlers *ws_handlers.WsHandlers,
) *Chat {
	return &Chat{clients: make(map[int64]map[*ClientChat]bool),
		register: make(chan *ClientChat),
		unregister: make(chan *ClientChat),
		action: make(chan *DataInfos),
		chatRepo: repo.Chat,
		userRepo: repo.User,
		chatHandler: wsHandlers.Chat,
	}
}

func JoinChat() gin.HandlerFunc  {
	return func(c *gin.Context) {
        userId, ok := c.Get("userId")
        if !ok {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "The userId is empty"})
            return
        }
		// TO CONTINUE
	}
}