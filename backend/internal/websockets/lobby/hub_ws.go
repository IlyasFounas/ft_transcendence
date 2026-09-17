package lobby_ws

import (
	"backend/internal/repositories"
	"backend/internal/websockets/handlers"
	"context"
	"encoding/json"
	"log"
	"time"
)

// duration of the pause between the end of a round and the start of the next one (if any), no pause if it was the last round.
const postRoundPause = 10 * time.Second

type kickPlayer struct {
	TargetId int64 `json:"target_id"`
}

type playerInfos struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar"`
}

// internalRoundEvent carries the game_id/round_nb for events the HUB ITSELF generates
// like : round_timeout firing, automatic start of the next round after the post-round pause
// -> these events never come from a client — they're sent on h.action from a timer callback (Client is nil in that case)
type internalRoundEvent struct {
	GameId  int64 `json:"game_id"`
	RoundNb int64 `json:"round_nb"`
}

// It will be sent to all clients when a round starts, with the round number and the absolute end time of the round
// (not a relative duration). Each client will calculate its own countdown using (RoundEndsAt - now), so all clients
// will converge to the same final instant regardless of their network latency at the time of reception.
type RoundStartedBroadcast struct {
	RoundNb     int64     `json:"round_nb"`
	RoundEndsAt time.Time `json:"round_ends_at"`
}

// It shares the results of the round to all players, even those who did not guess (they will see their own result as "not guessed").
// The results are sorted by player ID for consistency across clients.
// It is shared at the end of a round, with an indication of whether it was the last round or not :
// 	-> if it is not the last round, NextRoundAt is provided as an absolute timestamp for when the next round will start.
// It will be nil or omitted if IsLastRound is true.
type RoundResultBroadcast struct {
	Results     []ws_handlers.PlayerRoundResult `json:"results"`
	IsLastRound bool                            `json:"is_last_round"`
	NextRoundAt *time.Time                      `json:"next_round_at,omitempty"`
}

// encapsulation of the message sent to the client, with a type and a data payload interface{}
// that can be any struct (RoundStartedBroadcast, RoundResultBroadcast, etc.)
type WsMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// roundPhase describes what h.roundState[gameId] currently represents
// so a resyncing client gets the right kind of message (a running round vs the pause between rounds)
type roundPhase string

const (
	phasePlaying roundPhase = "playing" // a round is currently running: look at RoundEndsAt
	phasePause   roundPhase = "pause"   // between two rounds (or after the last one): look at LastResults/NextRoundAt
)

// to be checked when a client reconnects mid-game - refresh, dropped wifi, joining a game already in progress
type gameRoundState struct {
	Phase       roundPhase
	RoundNb     int64
	RoundEndsAt time.Time // valid when Phase == Playing
	NextRoundAt time.Time // valid when Phase == Pause and !IsLastRound
	IsLastRound bool
	LastResults []ws_handlers.PlayerRoundResult // results of the round that just ended, valid when Phase == phasePause
}

// DataInfos carries one action through the hubw
// it can comes from a real client or be generated internally by a timer callback
type DataInfos struct {
	Action  string
	Payload json.RawMessage
	// Client is nil for events generated internally by the hub
	Client *ClientWs
}

type HubWs struct {
	clients        map[int64]map[*ClientWs]bool
	register       chan *ClientWs
	unregister     chan *ClientWs
	action         chan *DataInfos
	roundTimers    map[int64]*time.Timer // means map[gameId]timer
	roundState     map[int64]*gameRoundState // means map[gameId]currentState, used to resync reconnecting clients
	lobbyRepo      *repositories.LobbyRepository
	userRepo       *repositories.UserRepository
	lobbyWsHandler *ws_handlers.LobbyWsHandler
	guessWsHandler *ws_handlers.GuessWsHandler
}

func NewHub(
	repo *repositories.Repositories,
	wsHandlers *ws_handlers.WsHandlers,
) *HubWs {
	return &HubWs{
		clients:        make(map[int64]map[*ClientWs]bool),
		register:       make(chan *ClientWs),
		unregister:     make(chan *ClientWs),
		action:         make(chan *DataInfos),
		roundTimers:    make(map[int64]*time.Timer),
		roundState:     make(map[int64]*gameRoundState),
		lobbyRepo:      repo.Lobby,
		userRepo:       repo.User,
		lobbyWsHandler: wsHandlers.Lobby,
		guessWsHandler: wsHandlers.Guess,
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
			ID:       user.ID,
			Username: user.Username,
			Avatar:   user.Avatar,
		})
	}
	h.broadcast(gameId, WsMessage{Type: "lobby_players", Data: usersLobby})
}

// broadcast serializes the payload once and sends it to all currently connected clients for gameId
// -> serialization means that the payload is converted to JSON format before being sent to the clients
// -> deserialization means that the backend will convert the JSON back to the original data structure when it receives it
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

// sendToClient serializes the payload once and sends it to ONE single client 
// unlike broadcast, it is used for resync messages that only concern the reconnecting client
func (h *HubWs) sendToClient(client *ClientWs, payload interface{}) {
	serializedData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	client.send <- serializedData
}

// resyncClient sends the newly reconnected client whatever the hub currently knows about the game's round state
func (h *HubWs) resyncClient(client *ClientWs) {
	state, ok := h.roundState[client.gameId]
	if !ok {
		// no round has started yet for this game (e.g. still in the lobby) : nothing to resync.
		return
	}
	switch state.Phase {
	case phasePlaying:
		h.sendToClient(client, WsMessage{
			Type: "round_started",
			Data: RoundStartedBroadcast{RoundNb: state.RoundNb, RoundEndsAt: state.RoundEndsAt.UTC()},
		})
	case phasePause:
		payload := RoundResultBroadcast{Results: state.LastResults, IsLastRound: state.IsLastRound}
		if !state.IsLastRound {
			nextRoundAt := state.NextRoundAt
			payload.NextRoundAt = &nextRoundAt
		}
		h.sendToClient(client, WsMessage{Type: "round_results", Data: payload})
	}
}

// this decides, from the hub's state, which round a "start_round" request from a client should actually restart from (but never from anything the client sent):
//   - no state at all: the game hasn't started, so it's round 1
//   - phasePlaying: a round is already running, so this is just a client resyncing
//     (StartRound/EnsureRound is idempotent for a round that already exists, which means it won't recreate it or extend its timer)
//   - phasePause: we're between two rounds, so it's the next one
func (h *HubWs) resolveRoundToStart(gameId int64) int64 {
	state, ok := h.roundState[gameId]
	if !ok {
		return 1
	}
	if state.Phase == phasePlaying || state.IsLastRound {
		return state.RoundNb
	}
	return state.RoundNb + 1
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
		// states' deletion (pause or playing) 
		if state, ok := h.roundState[client.gameId]; ok && state.IsLastRound {
			delete(h.roundState, client.gameId)
		}
	}
	close(client.send)
}

// Routes an action received (from a client or generated by the hub) to the appropriate handler.
// Handled sequentially (only one goroutine reads h.action in Run) : no race condition possible between
// two actions for the same gameId, even if they come from different clients.
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
		if data.Client != nil {
			// client != nil > real client request
			roundNb := h.resolveRoundToStart(data.Client.gameId)
			h.StartRound(ctx, data.Client.gameId, roundNb)
			return
		}
		// client == nil > this is the auto relaunch of the next round server-generated
		var evt internalRoundEvent
		if err := json.Unmarshal(data.Payload, &evt); err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		h.StartRound(ctx, evt.GameId, evt.RoundNb)

	case "round_timeout":
		if data.Client != nil {
			// round_timeout must only ever be generated internally by armRoundTimer
			log.Printf("warning : client %d tried to trigger round_timeout directly, ignored\n", data.Client.userId)
			return
		}
		var evt internalRoundEvent
		if err := json.Unmarshal(data.Payload, &evt); err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		h.EndRound(ctx, evt.GameId, evt.RoundNb)

	default:
		log.Printf("unknown action : %s\n", data.Action)
	}
}

// handleGuess persists the received guess, then ends the round if all active players in the game have guessed
func (h *HubWs) handleGuess(ctx context.Context, data *DataInfos) {
	createdGuess, ok := h.guessWsHandler.AddGuess(ctx, data.Payload, data.Client.userId, data.Client.gameId)
	if !ok {
		// TODO : renvoyer une erreur ciblée uniquement à data.Client
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
		h.EndRound(ctx, data.Client.gameId, createdGuess.RoundNb)
	}
}

// It arms the round timer for the given gameId and roundNb, and broadcasts the round_started event to all clients.
// The timer is set to the remaining duration of the round, calculated from the round's StartedAt timestamp and the round duration.
// If the round has already expired (remaining <= 0), it immediately ends the round instead of arming a negative timer.
// This ensures that clients who reconnect late or experience double triggers do not get more time than expected, and keeps all clients
// synchronized with the actual round start time.
func (h *HubWs) StartRound(ctx context.Context, gameId, roundNb int64) {
	// TODO vérif que le joueur est autorisé à démarrer un round
	round, err := h.guessWsHandler.EnsureRound(ctx, gameId, roundNb)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	duration, err := h.guessWsHandler.GetRoundDuration(ctx, gameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}

	roundEndsAt := round.StartedAt.Add(duration)
	remaining := time.Until(roundEndsAt)
	if remaining <= 0 {
		// round has expired (ex: late reco): ends it immediately
		h.EndRound(ctx, gameId, roundNb)
		return
	}

	h.roundState[gameId] = &gameRoundState{
		Phase:       phasePlaying,
		RoundNb:     roundNb,
		RoundEndsAt: roundEndsAt,
	}

	h.armRoundTimer(gameId, roundNb, remaining)
	h.broadcast(gameId, WsMessage{
		Type: "round_started",
		Data: RoundStartedBroadcast{RoundNb: roundNb, RoundEndsAt: roundEndsAt.UTC()},
	})
}

// this starts the countdown for the round's guess phase: when the timer expires,
// the round is ended even if not all players have guessed
func (h *HubWs) armRoundTimer(gameId, roundNb int64, duration time.Duration) {
	if existing, ok := h.roundTimers[gameId]; ok {
		existing.Stop()
	}
	h.roundTimers[gameId] = time.AfterFunc(duration, func() {
		payload, err := json.Marshal(internalRoundEvent{GameId: gameId, RoundNb: roundNb})
		if err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		h.action <- &DataInfos{Action: "round_timeout", Payload: payload}
	})
}

// Finalises the round roundNb: broadcasts the results of ALL active players (whether they guessed or not),
// then schedules a 10s pause before the next round — unless it was the last round.
func (h *HubWs) EndRound(ctx context.Context, gameId, roundNb int64) {
	if state, ok := h.roundState[gameId]; ok && state.Phase == phasePause && state.RoundNb == roundNb {
		// happens when the last guess of a round is processed at almost the same instant the round's timer fires
		return
	}
	if timer, ok := h.roundTimers[gameId]; ok {
		timer.Stop()
		delete(h.roundTimers, gameId)
	}

	roundID, err := h.guessWsHandler.ResolveRoundID(ctx, gameId, roundNb)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	activePlayers, err := h.lobbyRepo.AllPlayers(ctx, gameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	roundGuesses, err := h.guessWsHandler.GetRoundGuesses(ctx, roundID)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	results, err := h.guessWsHandler.BuildRoundResults(ctx, gameId, roundID, roundNb, activePlayers, roundGuesses)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	isLastRound, err := h.guessWsHandler.IsLastRound(ctx, gameId, roundNb)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}

	broadcastPayload := RoundResultBroadcast{Results: results, IsLastRound: isLastRound}
	var nextRoundAt time.Time
	if !isLastRound {
		// calculated once: this value is the reference. The AfterFunc below is based on the remaining time
		// until this same instant (time.Until), so the next round starts exactly when all clients have seen their countdown reach zero.
		nextRoundAt = time.Now().Add(postRoundPause)
		broadcastPayload.NextRoundAt = &nextRoundAt
	}
		h.roundState[gameId] = &gameRoundState{
		Phase:       phasePause,
		RoundNb:     roundNb,
		NextRoundAt: nextRoundAt,
		IsLastRound: isLastRound,
		LastResults: results,
	}

	h.broadcast(gameId, WsMessage{Type: "round_results", Data: broadcastPayload})

	if isLastRound {
		return
	}

	nextRoundNb := roundNb + 1
	time.AfterFunc(time.Until(nextRoundAt), func() {
		payload, err := json.Marshal(internalRoundEvent{GameId: gameId, RoundNb: nextRoundNb})
		if err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		h.action <- &DataInfos{Action: "start_round", Payload: payload}
	})
}

func (h *HubWs) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
			h.clientListBroadcast(client.gameId, ctx)
			// send this client (and only this client) the current round state, if any exists
			// so joining mid-game or reconnecting doesn't leave it blind until the next broadcast
			h.resyncClient(client) // TODO check if client reconnects only if it was already registered before (h.clients[gameId][client] == true) : if so, skip resyncClient because it already has the current state. Otherwise a reconnecting client will get two resyncs (one from the first registration, one from the second) and see the same round_started broadcast twice.
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