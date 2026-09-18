package lobby_ws

import (
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

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {return true},
}

type kickPlayer struct {
	TargetId int64 	 `json:"target_id"`
}

type playerInfos struct {
    ID       int64   `json:"id"`
    Username string  `json:"username"`
    Avatar   *string `json:"avatar"`
	Host	 bool	 `json:"host"`
}

type ClientHub struct {
    Connx     *client_ws.Connexion
    UserId  int64
    LobbyId int64
	GameId  int64
}

type DataInfos struct {
	Action	string
	Payload json.RawMessage
	Client	*ClientHub
}

type Hub struct {
	Clients map[int64]map[*ClientHub]bool
	Register chan *ClientHub
	Unregister chan *ClientHub
	Action chan *DataInfos
	LobbyRepo *repositories.LobbyRepository
	UserRepo *repositories.UserRepository
	LobbyWsHandler *ws_handlers.LobbyWsHandler
}


func NewHub(
	repo *repositories.Repositories,
	wsHandlers *ws_handlers.WsHandlers,
) *Hub {
	return &Hub{Clients: make(map[int64]map[*ClientHub]bool),
		Register: make(chan *ClientHub),
		Unregister: make(chan *ClientHub),
		Action: make(chan *DataInfos),
		LobbyRepo: repo.Lobby,
		UserRepo: repo.User,
		LobbyWsHandler: wsHandlers.Lobby,
	}
}

func NewClientHub(connx *client_ws.Connexion,
    userId int64, lobbyId int64, gameId int64) *ClientHub {
    return &ClientHub{Connx: connx,
    UserId: userId,
    LobbyId: lobbyId,
    GameId: gameId}
}

func (h *Hub) broadcast(gameId int64, payload interface{}) {
    serializedData, err := json.Marshal(payload)
    if err != nil {
        log.Printf("error : %s\n", err)
        return
    }
    for client := range h.Clients[gameId] {
        client.Connx.Send <- serializedData
    }
}

func (h *Hub) broadcastClient(gameId int64, userId int64, payload interface{}) {
    serializedData, err := json.Marshal(payload)
    if err != nil {
        log.Printf("error : %s\n", err)
        return
    }
	for client := range h.Clients[gameId] {
		if client.UserId == userId {
			client.Connx.Send <- serializedData
		}	
	}
}

func (h *Hub) clientLobbyBroadcast(gameId int64, ctx context.Context) {
	lobbies, err := h.LobbyRepo.AllPlayers(ctx, gameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return 
	}
	usersLobby := make([]playerInfos, 0, len(lobbies))
	for _, lobby := range lobbies {
		user, err := h.UserRepo.User(ctx, lobby.UserId)
		if err != nil {
			log.Printf("error : %s\n", err)
			return 
		}
		usersLobby = append(usersLobby, playerInfos{
			ID:       user.ID,
			Username: user.Username,
			Avatar:   user.Avatar,
			Host:     lobby.IsCreator,
		})
	}
	h.broadcast(gameId, usersLobby)
}


func (h *Hub) registerClient(client *ClientHub) {
	if _, ok := h.Clients[client.GameId]; !ok {
		h.Clients[client.GameId] = make(map[*ClientHub]bool)
	}
	h.Clients[client.GameId][client] = true
}

func (h *Hub) unregisterClient(client *ClientHub) {
	if _, ok := h.Clients[client.GameId]; ok {
		delete(h.Clients[client.GameId], client)
	}
	if len(h.Clients[client.GameId]) == 0 {
		delete(h.Clients, client.GameId)
	}
}

func (h *Hub) handleAction(ctx context.Context,
	serializedData *DataInfos) {
	switch serializedData.Action {
		case "leave":
			serializedData.Client.Connx.Conn.Close()
		case "kick_player":
			var targetClient kickPlayer
			err := json.Unmarshal(serializedData.Payload, &targetClient)
			if err != nil {
				log.Printf("error : %s\n", err)
				return 
			}
			if !h.LobbyWsHandler.RemovePlayer(ctx,
				serializedData.Client.UserId, targetClient.TargetId) {
					return 
			}
			for client := range h.Clients[serializedData.Client.GameId] {
				if client.UserId == targetClient.TargetId {
					h.clientLobbyBroadcast(serializedData.Client.GameId, ctx)
					client.Connx.Conn.Close()
					return 
				}
			}
			h.clientLobbyBroadcast(serializedData.Client.GameId, ctx)
		
		default:
			log.Printf("unknown action : %s\n", serializedData.Action)
	}
}

func JoinLobby(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "The userId is empty"})
			return
		}        
		// Check if the player is in the right lobby
        lobby, err := hub.LobbyRepo.ActualLobby(c.Request.Context(), userId.(int64))
        if err != nil || lobby == nil {
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "Can't find the lobby"})
            return
        }
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("error : %s\n", err)
			return 
		}
		connx := client_ws.NewConnexion(conn)
		clientHub := NewClientHub(connx, userId.(int64),
			lobby.ID, lobby.GameId)
		hub.Register <- clientHub
		go clientHub.Connx.WritePump()
		actionClient := func(data []byte) {
			var dataInfos DataInfos
			err := json.Unmarshal(data, &dataInfos)
			if err != nil {
				log.Printf("error : %s\n", err)
				return
			}
			dataInfos.Client = clientHub
			hub.Action <- &dataInfos 
		}
		unregisterClient := func() {
			hub.Unregister <- clientHub
		}
		clientHub.Connx.ReadPump(actionClient, unregisterClient)
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <- h.Register:
			h.registerClient(client)
			h.clientLobbyBroadcast(client.GameId, ctx)
			data, ok := h.LobbyWsHandler.SendInfos(ctx, client.GameId)
			if !ok {
				continue
			}
			h.broadcastClient(client.GameId, client.UserId, data)
		case client := <- h.Unregister:
			h.unregisterClient(client)
			if !h.LobbyWsHandler.LeaveLobby(ctx, client.UserId) {
				log.Println("error : the client can't leave the lobby in the database")
				continue 
			}
			h.clientLobbyBroadcast(client.GameId, ctx)
		case serializedData := <- h.Action:
			h.handleAction(ctx, serializedData)
		}
	}
}