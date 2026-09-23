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

// logAndFail logs err (using this file's shared "error : %s" format) and reports whether the
// caller should abort. Replaces the "if err != nil { log.Printf(...); return false }" block that
// used to be repeated, verbatim, after almost every repository call in this file.
func logAndFail(err error) bool {
	if err != nil {
		log.Printf("error : %s\n", err)
		return true
	}
	return false
}

func (h *LobbyWsHandler) LeaveLobby(ctx context.Context,
	userId int64) bool {
	lobby, err := h.lobbyRepo.ActualLobby(ctx, userId)
	if err != nil || lobby == nil {
		log.Printf("error : %s\n", err)
		return false
	}
	if logAndFail(h.lobbyRepo.SetLeftAt(ctx, lobby.ID)) {
		return false
	}
	lobbies, err := h.lobbyRepo.AllPlayers(ctx, lobby.GameId)
	if logAndFail(err) {
		return false
	}
	if logAndFail(h.lobbyRepo.SetJoinCode(ctx, lobby.ID, nil)) {
		return false
	}
	if logAndFail(h.lobbyRepo.SetIsCreator(ctx, lobby.ID, false)) {
		return false
	}
	for _, newLobby := range lobbies {
		if newLobby.UserId != userId {
			if logAndFail(h.lobbyRepo.SetIsCreator(ctx, newLobby.ID, true)) {
				return false
			}
			if logAndFail(h.lobbyRepo.SetJoinCode(ctx, newLobby.ID, lobby.JoinCode)) {
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
	if logAndFail(err) {
		return false
	}
	lobby, err := h.lobbyRepo.ActualLobby(ctx, user.ID)
	if logAndFail(err) {
		return false
	}
	isCreator, err := h.lobbyRepo.IsCreator(ctx, request_id, lobby.GameId)
	if logAndFail(err) {
		return false
	}
	if !isCreator {
		log.Println("error : you're not allowed to remove the player")
		return false
	}
	if logAndFail(h.lobbyRepo.SetLeftAt(ctx, lobby.ID)) {
		return false
	}
	if logAndFail(h.lobbyRepo.SetIsCreator(ctx, lobby.ID, false)) {
		return false
	}
	return true
}

// fetchInt64 is the shared body behind Timer and RoundTotal: call fn, log and report failure if
// it errors, otherwise return its value. Timer and RoundTotal stay separate named methods (called
// by name from SendInfos and elsewhere), but no longer duplicate the same four-line body twice.
func (h *LobbyWsHandler) fetchInt64(ctx context.Context, fn func(context.Context, int64) (int64, error), gameId int64) (int64, bool) {
	value, err := fn(ctx, gameId)
	if logAndFail(err) {
		return 0, false
	}
	return value, true
}

func (h *LobbyWsHandler) Timer(ctx context.Context, gameId int64) (int64, bool) {
	return h.fetchInt64(ctx, h.gameRepo.Timer, gameId)
}

func (h *LobbyWsHandler) RoundTotal(ctx context.Context, gameId int64) (int64, bool) {
	return h.fetchInt64(ctx, h.gameRepo.RoundTotal, gameId)
}

func (h *LobbyWsHandler) SendInfos(ctx context.Context, gameId int64) (interface{}, bool) {
	timer, ok := h.Timer(ctx, gameId)
	if !ok {
		return nil, false
	}
	joinCode, err := h.lobbyRepo.JoinCode(ctx, gameId)
	if logAndFail(err) {
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

// UpdateStartAt records when the game actually started. Only ever called once, right after
// RoundWsHandler.RequestGameStart successfully schedules round 1's pre-round pause.
func (h *LobbyWsHandler) UpdateStartAt(ctx context.Context, LobbyId int64) error {
	if err := h.lobbyRepo.SetStartAt(ctx, LobbyId); err != nil {
		log.Printf("error : %s\n", err)
		return err
	}
	return nil
}