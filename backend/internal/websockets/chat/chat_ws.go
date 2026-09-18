package chat_ws

import (
	"backend/internal/models/friends"
	"backend/internal/repositories"
	"backend/internal/websockets/client"
	"backend/internal/websockets/handlers"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Chat struct {
	Clients     map[int64]map[*ClientChat]bool
	Register    chan *ClientChat
	Unregister  chan *ClientChat
	Action      chan *DataInfos
	ChatRepo    *repositories.ChatRepository
	ChatHandler *ws_handlers.ChatHandler
	UserRepo    *repositories.UserRepository
}

type ClientInfo struct {
	Id   		int64
	Username	string
	Avatar 		string
	Online		bool
}

type ClientChat struct {
	Connx    *client_ws.Connexion
	UserId   int64
	TargetId int64
}

type DataInfos struct {
	Action  string
	Payload json.RawMessage
	Client  *ClientChat
}

type ChatFriend struct {
	ID			int64   	`json:"id"`
	UserID     	int64     	`json:"user_id"`
	Name		string    	`json:"name"`
	Avatar    	*string   	`json:"avatar"`
	Online		bool		`json:"online"`
}

func NewChat(
	repo *repositories.Repositories,
	wsHandlers *ws_handlers.WsHandlers,
) *Chat {
	return &Chat{Clients: make(map[int64]map[*ClientChat]bool),
		Register:    make(chan *ClientChat),
		Unregister:  make(chan *ClientChat),
		Action:      make(chan *DataInfos),
		ChatRepo:    repo.Chat,
		ChatHandler: wsHandlers.Chat,
		UserRepo:    repo.User,
	}
}

func NewClientChat(
	connx *client_ws.Connexion,
	userId int64,
) *ClientChat {
	return &ClientChat{Connx: connx,
		UserId:   userId,
	}
}

func (c *Chat) IsOnline(userId int64) bool {
	if _, ok := c.Clients[userId]; !ok {
		return false
	}
	return true
}

func (c *Chat) sendToClient(userId int64, friends []*friends.FriendResponse) {
	friendsList := make([]ChatFriend, 0, len(friends))
	for _, friend := range friends {
		friendsList = append(friendsList, ChatFriend{
			ID : friend.ID,
			UserID : friend.UserID,
			Name : friend.Name,
			Avatar : friend.Avatar,
			Online : c.IsOnline(friend.UserID),
		})
	}
	data, err := json.Marshal(friendsList)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
    for client := range c.Clients[userId] {
		client.Connx.Send <- data
	}
}

func (c *Chat) registerClient(client *ClientChat) {
	if _, ok := c.Clients[client.UserId]; !ok {
		c.Clients[client.UserId] = make(map[*ClientChat]bool)
	}
	c.Clients[client.UserId][client] = true
}

func (c *Chat) unregisterClient(client *ClientChat) {
	if _, ok := c.Clients[client.UserId]; ok {
		delete(c.Clients[client.UserId], client)
	}
	if len(c.Clients[client.UserId]) == 0 {
		delete(c.Clients, client.UserId)
	}
}

func (c *Chat) handleAction(serializedData *DataInfos) {
	switch serializedData.Action {
	case "open_chat":
		// TODO open the chat and give the 20 last messages
	case "send":
		// TODO send the message to the friend
	case "receive_messages":
		// TODO give the 20 more last messages
	default:
		log.Printf("unknown action : %s\n", serializedData.Action)
	}
}

func JoinChat(chat *Chat) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		connx := client_ws.NewConnexion(conn)
		clientChat := NewClientChat(connx, userId.(int64))
		chat.Register <- clientChat
		go clientChat.Connx.WritePump()
		actionClient := func(data []byte) {
			var dataInfos DataInfos
			err := json.Unmarshal(data, &dataInfos)
			if err != nil {
				log.Printf("error : %s\n", err)
				return
			}			
			dataInfos.Client = clientChat
			chat.Action <- &dataInfos
		}
		unregisterClient := func() {
			chat.Unregister <- clientChat
		}
		clientChat.Connx.ReadPump(actionClient,
			unregisterClient)
	}
}

func (c *Chat) Run(ctx context.Context) {
	for {
		select {
		case client := <-c.Register:
			c.registerClient(client)
			friends, ok := c.ChatHandler.FriendList(ctx, client.UserId)
			if !ok {
				continue 
			}
			c.sendToClient(client.UserId, friends)
		case client := <-c.Unregister:
			c.unregisterClient(client)
		case data := <-c.Action:
			c.handleAction(data)
		}
	}
}
