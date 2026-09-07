package lobby_ws

// import (
// )

type HubWs struct {
	clients map[int64]map[*ClientWs]bool
	register chan *ClientWs
	unregister chan *ClientWs
}

func NewHub() *HubWs {
	return &HubWs{clients: make(map[int64]map[*ClientWs]bool),
		register: make(chan *ClientWs),
		unregister: make(chan *ClientWs)}
}

func (h *HubWs) registerClient(client *ClientWs) {
	if _, ok := h.clients[client.lobbyId]; !ok {
		h.clients[client.lobbyId] = make(map[*ClientWs]bool)
	}
	h.clients[client.lobbyId][client] = true
}

func (h *HubWs) unregisterClient(client *ClientWs) {
	if _, ok := h.clients[client.lobbyId]; ok {
		delete(h.clients[client.lobbyId], client)
	}
}

func (h *HubWs) Run() {
	for {
		select {
		case client := <- h.register:
			h.registerClient(client)
		case client := <- h.unregister:
			h.unregisterClient(client)
		}
	}
}