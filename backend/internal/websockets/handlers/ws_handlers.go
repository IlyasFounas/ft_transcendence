package ws_handlers

import (
	"backend/internal/repositories"
)

type WsHandlers struct {
	Lobby	*LobbyWsHandler
	Chat	*ChatHandler
}

func NewWsHandlers(repo *repositories.Repositories) *WsHandlers {
	return &WsHandlers{Lobby: NewLobbyWsHandler(repo.Lobby, repo.Game, repo.User),
		Chat: NewChatHandler(repo.Chat, repo.Friend, repo.User)}
}