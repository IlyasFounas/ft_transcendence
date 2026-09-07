package lobby_ws

import (
	// "backend/internal/websockets/lobby"
)

type HubWs struct {
	clients map[int64]map[*ClientWs]bool
	register chan *ClientWs
	unregister chan *ClientWs
}

func NewHub() *HubWs {
	return &HubWs{clients: make(map[int64]map[*ClientWs]bool)}
}

func (h *HubWs) registerClient(client *ClientWs) {
	if _, ok := h.clients[client.lobbyId]; !ok {
		h.clients[client.lobbyId] = make(map[*ClientWs]bool)
	}
	h.clients[client.lobbyId][client] = true
}

func (h *HubWs) unregisterClient() {
	if _, ok := h.clients[client.lobbyId]; ok {
		h.clients[client.lobbyId] = make(map[*ClientWs]bool)
		delete()
	}
}

func (h *HubWs) Run() {
	for {
		select {
		case client := <- h.register:
			registerClient(client)
		case client := <- h.unregister:
			// do smth
		}
	}
}