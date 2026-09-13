package ws_handlers

import (
	"backend/internal/repositories"
)

type WsHandlers struct {
	Lobby	*LobbyWsHandler
}

func NewWsHandlers(repo *repositories.Repositories) *WsHandlers {
	return &WsHandlers{Lobby: NewLobbyWsHandler(repo.Lobby, repo.User)}
}