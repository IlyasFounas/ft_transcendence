package lobby_ws

import (
	"backend/internal/models/games"
	"backend/internal/models/guesses"
	"backend/internal/repositories"
	"backend/internal/websockets/handlers"
	"context"
	"encoding/json"
	"log"
	"time"
)

type kickPlayer struct {
	TargetId	int64 `json:"target_id"`
}

type playerInfos struct {
	ID			int64   `json:"id"`
	Username	string  `json:"username"`
	Avatar		*string `json:"avatar"`
}

type startRoundPayload struct {
	RoundID	int64 `json:"round_id"`
	RoundNb	int64 `json:"round_nb"`
}

type roundTimeoutPayload struct {
	RoundID	int64 `json:"round_id"`
	RoundNb	int64 `json:"round_nb"`
}

type roundStartedPayload struct {
	RoundID		int64  `json:"round_id"`
	RoundEndsAt	string `json:"round_ends_at"` // RFC3339 UTC
}

type DataInfos struct {
	Action	string
	Payload	json.RawMessage
	Client	*ClientWs
	GameId	int64
}

func (d *DataInfos) resolveGameId() int64 {
	if d.Client != nil {
		return d.Client.gameId
	}
	return d.GameId
}

type HubWs struct {
	clients			map[int64]map[*ClientWs]bool
	register		chan *ClientWs
	unregister		chan *ClientWs
	action			chan *DataInfos
	roundTimers		map[int64]*time.Timer // clé = gameId, uniquement touché depuis Run()
	lobbyRepo		*repositories.LobbyRepository
	userRepo		*repositories.UserRepository
	gameRepo		*repositories.GameRepository
	lobbyWsHandler	*ws_handlers.LobbyWsHandler
	guessWsHandler	*ws_handlers.GuessWsHandler
}

func NewHub(
	repo *repositories.Repositories,
	wsHandlers *ws_handlers.WsHandlers,
) *HubWs {
	return &HubWs{
		clients:		make(map[int64]map[*ClientWs]bool),
		register:		make(chan *ClientWs),
		unregister:		make(chan *ClientWs),
		action:			make(chan *DataInfos),
		roundTimers:	make(map[int64]*time.Timer),
		lobbyRepo:		repo.Lobby,
		userRepo:		repo.User,
		gameRepo:		repo.Game,
		lobbyWsHandler:	wsHandlers.Lobby,
		guessWsHandler:	wsHandlers.Guess,
	}
}

func (h *HubWs) clientListBroadcast(gameId int64, ctx context.Context) {
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
			ID:			user.ID,
			Username:	user.Username,
			Avatar:		user.Avatar,
		})
	}
	h.broadcast(gameId, usersLobby)
}

// broadcast sérialise payload une seule fois et l'envoie à tous les
// clients actuellement connectés pour gameId.
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
	close(client.send)
}

func (h *HubWs) handleAction(ctx context.Context, data *DataInfos) {
	switch data.Action {
	case "leave":
		data.Client.conn.Close()

	case "kick_player":
		var targetClient kickPlayer
		if err := json.Unmarshal(data.Payload, &targetClient); err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		if !h.lobbyWsHandler.RemovePlayer(ctx, data.Client.userId, targetClient.TargetId) {
			return
		}
		for client := range h.clients[data.Client.gameId] {
			if client.userId == targetClient.TargetId {
				h.clientListBroadcast(data.Client.gameId, ctx)
				client.conn.Close()
				return
			}
		}
		h.clientListBroadcast(data.Client.gameId, ctx)

	case "guess":
		h.handleGuess(ctx, data)

	case "start_round":
		var payload startRoundPayload
		if err := json.Unmarshal(data.Payload, &payload); err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		gameId := data.resolveGameId()

		timerSeconds, err := h.gameRepo.Timer(ctx, gameId)
		if err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		duration := time.Duration(timerSeconds) * time.Second

		h.startRoundTimer(gameId, payload.RoundID, payload.RoundNb, duration)
		h.broadcast(gameId, roundStartedPayload{
			RoundID:		payload.RoundID,
			RoundEndsAt:	time.Now().UTC().Add(duration).Format(time.RFC3339),
		})

	case "round_timeout":
		var payload roundTimeoutPayload
		if err := json.Unmarshal(data.Payload, &payload); err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		gameId := data.resolveGameId()
		delete(h.roundTimers, gameId)

		activePlayers, err := h.lobbyRepo.AllPlayers(ctx, gameId)
		if err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		roundGuesses, err := h.guessWsHandler.GetRoundGuesses(ctx, payload.RoundID)
		if err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		h.finalizeRound(ctx, gameId, payload.RoundID, payload.RoundNb, activePlayers, roundGuesses)

	default:
		log.Printf("unknown action : %s\n", data.Action)
	}
}


func (h *HubWs) handleGuess(ctx context.Context, data *DataInfos) {
	createdGuess, ok := h.guessWsHandler.AddGuess(ctx, data.Payload, data.Client.userId, data.Client.gameId)
	if !ok {
		return
	}

	activePlayers, err := h.lobbyRepo.AllPlayers(ctx, data.Client.gameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}

	roundGuesses, err := h.guessWsHandler.GetRoundGuesses(ctx, createdGuess.RoundID)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}

	if len(roundGuesses) >= len(activePlayers) {
		h.finalizeRound(ctx, data.Client.gameId, createdGuess.RoundID, createdGuess.RoundNb, activePlayers, roundGuesses)
	}
}


func (h *HubWs) finalizeRound(ctx context.Context, gameId, roundId, roundNb int64,
	activePlayers []*games.Lobby, roundGuesses []*guesses.Guess) {
	if timer, ok := h.roundTimers[gameId]; ok {
		timer.Stop()
		delete(h.roundTimers, gameId)
	}

	results, err := h.guessWsHandler.BuildRoundResults(ctx, gameId, roundId, roundNb, activePlayers, roundGuesses)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	h.broadcast(gameId, results)
}

func (h *HubWs) startRoundTimer(gameId, roundId, roundNb int64, duration time.Duration) {
	if existing, ok := h.roundTimers[gameId]; ok {
		existing.Stop()
	}
	h.roundTimers[gameId] = time.AfterFunc(duration, func() {
		payload, _ := json.Marshal(roundTimeoutPayload{RoundID: roundId, RoundNb: roundNb})
		h.action <- &DataInfos{Action: "round_timeout", Payload: payload, GameId: gameId}
	})
}

func (h *HubWs) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
			h.clientListBroadcast(client.gameId, ctx)
		case client := <-h.unregister:
			h.unregisterClient(client)
			if !h.lobbyWsHandler.LeaveLobby(ctx, client.userId) {
				log.Println("error : the client can't leave the lobby in the database")
				continue
			}
			h.clientListBroadcast(client.gameId, ctx)
		case data := <-h.action:
			h.handleAction(ctx, data)
		}
	}
}