package ws_handlers

import (
	"log"
	"context"
	"backend/internal/repositories"
)

type LobbyWsHandler struct {
	lobbyRepo *repositories.LobbyRepository
	userRepo *repositories.UserRepository
}

func NewLobbyWsHandler(lobbyRepo *repositories.LobbyRepository, 
	userRepo *repositories.UserRepository) *LobbyWsHandler {
	return &LobbyWsHandler{lobbyRepo: lobbyRepo,
		userRepo: userRepo}
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

func (h *LobbyWsHandler) RemovePlayer(ctx context.Context,
	request_id int64, target_id int64) bool {
	if request_id == target_id {
		log.Println("error : you can't kick yourself out of the lobby")
		return false	
	}
	user, err := h.userRepo.User(ctx, target_id)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	lobby, err := h.lobbyRepo.ActualLobby(ctx, user.ID)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	isCreator, err := h.lobbyRepo.IsCreator(ctx,
		request_id, lobby.GameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	if !isCreator {
		log.Println("error : you're not allowed to remove the player")
		return false
	}
	err = h.lobbyRepo.SetLeftAt(ctx, lobby.ID)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	return true
}