package ws_handlers

import (
	"backend/internal/repositories"
	"context"
	"database/sql"
	"errors"
	"log"
)

type LobbyWsHandler struct {
	lobbyRepo *repositories.LobbyRepository
	gameRepo  *repositories.GameRepository
	userRepo  *repositories.UserRepository
}

type lobbyInfos struct {
	JoinCode   *string `json:"join_code"`
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
	isCreator, err := h.lobbyRepo.IsCreator(ctx,
		userId, lobby.GameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	if isCreator {
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
	if lobby == nil {
		// the target is in no lobby (already kicked or gone)
		log.Println("error : the player to kick is not in this lobby")
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
	// the kicked player won't be able to join this game again
	err = h.lobbyRepo.Kick(ctx, lobby.ID)
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
	// the join code is removed when the game starts: a player reconnecting to a started game still gets the other infos
	joinCode, err := h.lobbyRepo.JoinCode(ctx, gameId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("error : %s\n", err)
		return nil, false
	}
	roundTotal, ok := h.RoundTotal(ctx, gameId)
	if !ok {
		return nil, false
	}
	data := lobbyInfos{
		JoinCode:   joinCode,
		RoundTotal: roundTotal,
		Timer:      timer,
	}
	return data, true
}

// Tells whether userId is the current host of the game gameId
func (h *LobbyWsHandler) IsHost(ctx context.Context, userId int64, gameId int64) bool {
	lobby, err := h.lobbyRepo.ActualLobby(ctx, userId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	return lobby != nil && lobby.GameId == gameId && lobby.IsCreator
}

// When the game starts and when it ends: the join code is removed (nobody can join anymore)
// and the players who are in the lobby without being connected to the WS (not in connected) are removed from it
func (h *LobbyWsHandler) CloseLobby(ctx context.Context, gameId int64, connected map[int64]bool) bool {
	players, err := h.lobbyRepo.AllPlayers(ctx, gameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	for _, player := range players {
		if player.JoinCode != nil {
			if err := h.lobbyRepo.SetJoinCode(ctx, player.ID, nil); err != nil {
				log.Printf("error : %s\n", err)
				return false
			}
		}
		if !connected[player.UserId] {
			if err := h.lobbyRepo.SetLeftAt(ctx, player.ID); err != nil {
				log.Printf("error : %s\n", err)
				return false
			}
		}
	}
	return true
}

// Removes userId from the game gameId once it is finished
// Nothing is handed over to another player (host, join code): the game is over, nobody needs them anymore
func (h *LobbyWsHandler) LeaveFinishedGame(ctx context.Context, userId int64, gameId int64) bool {
	lobby, err := h.lobbyRepo.ActualLobbyGameId(ctx, userId, gameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	if lobby == nil {
		// already out of the game (another tab left first, or removed at the end of the game)
		return true
	}
	if err := h.lobbyRepo.SetLeftAt(ctx, lobby.ID); err != nil {
		log.Printf("error : %s\n", err)
		return false
	}
	return true
}