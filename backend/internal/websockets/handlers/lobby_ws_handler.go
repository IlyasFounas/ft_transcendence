package ws_handlers

import (
	"log"
	"context"
	"backend/internal/repositories"
)

type LobbyWsHandler struct {
	lobbyRepo *repositories.LobbyRepository
}

func NewLobbyWsHandler(lobbyRepo *repositories.LobbyRepository) *LobbyWsHandler {
	return &LobbyWsHandler{lobbyRepo: lobbyRepo}
}

func (h *LobbyWsHandler) LeaveLobby(ctx context.Context,
	userId int64) bool {
	lobby, err := h.lobbyRepo.ActualLobby(ctx,
			userId);
	if err != nil || lobby == nil {
		log.Printf("error : %s\n", err)
		return false
	}
	err = h.lobbyRepo.SetLeftAt(ctx, lobby.ID)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	lobbies, err := h.lobbyRepo.AllPlayers(ctx, lobby.GameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	err = h.lobbyRepo.SetJoinCode(ctx,
		lobby.ID, nil)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	err = h.lobbyRepo.SetIsCreator(ctx,
		lobby.ID, false)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	for _, newLobby := range lobbies {
		if newLobby.UserId != userId {
			err = h.lobbyRepo.SetIsCreator(ctx,
				newLobby.ID, true)
			if err != nil {
				log.Printf("error : %s\n", err)
				return false
			}
			err = h.lobbyRepo.SetJoinCode(ctx,
				newLobby.ID, lobby.JoinCode)
			if err != nil {
				log.Printf("error : %s\n", err)
				return false
			}
			break
		}
	}
	return true
}
