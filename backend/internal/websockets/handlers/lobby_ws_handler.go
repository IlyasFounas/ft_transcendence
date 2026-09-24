package ws_handlers

import (
	"backend/internal/repositories"
	"context"
	"log"
)

type LobbyWsHandler struct {
	lobbyRepo *repositories.LobbyRepository
	gameRepo  *repositories.GameRepository
	userRepo  *repositories.UserRepository
}

type lobbyInfos struct {
	JoinCode   string `json:"join_code"`
	RoundTotal int64  `json:"round_total"`
	Timer      int64  `json:"timer"`
}

func NewLobbyWsHandler(lobbyRepo *repositories.LobbyRepository,
	gameRepo *repositories.GameRepository,
	userRepo *repositories.UserRepository) *LobbyWsHandler {
	return &LobbyWsHandler{lobbyRepo: lobbyRepo,
		gameRepo: gameRepo,
		userRepo: userRepo}
}

func (h *LobbyWsHandler) LeaveLobby(ctx context.Context,
	userId int64) bool {
	lobby, err := h.lobbyRepo.ActualLobby(ctx,
		userId)
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
	err = h.lobbyRepo.SetIsCreator(ctx,
		lobby.ID, false)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	return true
}

func (h *LobbyWsHandler) Timer(ctx context.Context, gameId int64) (int64, bool) {
	timer, err := h.gameRepo.Timer(ctx, gameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return 0, false
	}
	return timer, true
}

func (h *LobbyWsHandler) RoundTotal(ctx context.Context, gameId int64) (int64, bool) {
	roundTotal, err := h.gameRepo.RoundTotal(ctx, gameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return 0, false
	}
	return roundTotal, true
}

func (h *LobbyWsHandler) SendInfos(ctx context.Context, gameId int64) (interface{}, bool) {
	timer, ok := h.Timer(ctx, gameId)
	if !ok {
		return nil, false
	}
	joinCode, err := h.lobbyRepo.JoinCode(ctx, gameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, false
	}
	roundTotal, ok := h.RoundTotal(ctx, gameId)
	if !ok {
		return nil, false
	}
	data := lobbyInfos{
		JoinCode:   *joinCode,
		RoundTotal: roundTotal,
		Timer:      timer,
	}
	return data, true
}

func (h *LobbyWsHandler) UpdateStartAt(ctx context.Context, LobbyId int64, RoundNb int64) error {
	if RoundNb == 1 {
		err := h.lobbyRepo.SetStartAt(ctx, LobbyId)
		if err != nil {
			log.Printf("error : %s\n", err)
			return err
		}
	}
	return nil
}