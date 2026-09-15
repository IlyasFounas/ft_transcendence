package lobby_ws

import (
	"backend/internal/repositories"
	"backend/internal/websockets/handlers"
	"context"
	"encoding/json"
	"log"
)

type kickPlayer struct {
	TargetId int64 	 `json:"target_id"`
}

type playerInfos struct {
    ID       int64   `json:"id"`
    Username string  `json:"username"`
    Avatar   *string `json:"avatar"`
	Host	 bool	 `json:"host"`
}

type DataInfos struct {
	Action	string
	Payload json.RawMessage
	Client	*ClientWs
}

type HubWs struct {
	clients map[int64]map[*ClientWs]bool
	register chan *ClientWs
	unregister chan *ClientWs
	action chan *DataInfos
	lobbyRepo *repositories.LobbyRepository
	userRepo *repositories.UserRepository
	lobbyWsHandler *ws_handlers.LobbyWsHandler
}

func NewHub(
	repo *repositories.Repositories,
	wsHandlers *ws_handlers.WsHandlers,
) *HubWs {
	return &HubWs{clients: make(map[int64]map[*ClientWs]bool),
		register: make(chan *ClientWs),
		unregister: make(chan *ClientWs),
		action: make(chan *DataInfos),
		lobbyRepo: repo.Lobby,
		userRepo: repo.User,
		lobbyWsHandler: wsHandlers.Lobby,
	}
}

func (h *HubWs) broadcast(gameId int64, payload interface{}) {
    serializedData, err := json.Marshal(payload)
    if err != nil {
        log.Printf("error : %s\n", err)
        return
    }
    for client := range h.clients[gameId] {
        client.send <- serializedData
    }
}

func (h *HubWs) broadcastClient(gameId int64, userId int64, payload interface{}) {
    serializedData, err := json.Marshal(payload)
    if err != nil {
        log.Printf("error : %s\n", err)
        return
    }
	for client := range h.clients[gameId] {
		if client.userId == userId {
			client.send <- serializedData
		}	
	}
}

func (h *HubWs) clientLobbyBroadcast(gameId int64, ctx context.Context) {
	lobbies, err := h.lobbyRepo.AllPlayers(ctx, gameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return 
	}
	usersLobby := make([]playerInfos, 0, len(lobbies))
	for _, lobby := range lobbies {
		user, err := h.userRepo.User(ctx, lobby.UserId)
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


func (h *HubWs) registerClient(client *ClientWs) {
	if _, ok := h.clients[client.gameId]; !ok {
		h.clients[client.gameId] = make(map[*ClientWs]bool)
	}
	h.clients[client.gameId][client] = true
}

func (h *HubWs) unregisterClient(client *ClientWs) {
	if _, ok := h.clients[client.gameId]; ok {
		delete(h.clients[client.gameId], client)
	}
	if len(h.clients[client.gameId]) == 0 {
		delete(h.clients, client.gameId)
	}
}

func (h *HubWs) handleAction(ctx context.Context,
	serializedData *DataInfos) {
	switch serializedData.Action {
		case "leave":
			serializedData.Client.conn.Close()
		case "kick_player":
			var targetClient kickPlayer
			err := json.Unmarshal(serializedData.Payload, &targetClient)
			if err != nil {
				log.Printf("error : %s\n", err)
				return 
			}
			if !h.lobbyWsHandler.RemovePlayer(ctx,
				serializedData.Client.userId, targetClient.TargetId) {
					return 
			}
			for client := range h.clients[serializedData.Client.gameId] {
				if client.userId == targetClient.TargetId {
					h.clientLobbyBroadcast(serializedData.Client.gameId, ctx)
					client.conn.Close()
					return 
				}
			}
			h.clientLobbyBroadcast(serializedData.Client.gameId, ctx)
		
		default:
			log.Printf("unknown action : %s\n", serializedData.Action)
	}
}

func (h *HubWs) Run(ctx context.Context) {
	for {
		select {
		case client := <- h.register:
			h.registerClient(client)
			h.clientLobbyBroadcast(client.gameId, ctx)
			data, ok := h.lobbyWsHandler.SendInfos(ctx, client.gameId)
			if !ok {
				continue
			}
			h.broadcastClient(client.gameId, client.userId, data)
		case client := <- h.unregister:
			h.unregisterClient(client)
			if !h.lobbyWsHandler.LeaveLobby(ctx, client.userId) {
				log.Println("error : the client can't leave the lobby in the database")
				continue 
			}
			h.clientLobbyBroadcast(client.gameId, ctx)
		case serializedData := <- h.action:
			h.handleAction(ctx, serializedData)
		}
	}
}