package lobby_ws

import (
	"backend/internal/repositories"
	"backend/internal/websockets/client"
	"backend/internal/websockets/handlers"
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {return true},
}

// multiplayer only: pause before EVERY round, the first one included (to get ready before round 1, then to look at
// the results between two rounds). No pause after the last round, and none at all in solo (see handleStartGame).
const preRoundPause = 10 * time.Second

type kickPlayer struct {
	TargetId int64 `json:"target_id"`
}

type playerInfos struct {
	ID			int64	`json:"id"`
	Username	string	`json:"username"`
	Avatar		*string	`json:"avatar"`
	Host		bool	`json:"host"`
}

type ClientHub struct {
    Connx     *client_ws.Connexion
    UserId  int64
    LobbyId int64
	GameId  int64
}

// internalRoundEvent carries the game_id/round_nb for events the HUB ITSELF generates
// like : round_timeout firing, automatic start of the next round at the end of a multiplayer pause
// -> these events never come from a client — they're sent on h.action from a timer callback (Client is nil in that case)
type internalRoundEvent struct {
	GameId  int64 `json:"game_id"`
	RoundNb int64 `json:"round_nb"`
}

// Sent to all clients when the host starts the game, and resent to a client that (re)connects during the countdown before round 1.

type GameStartedBroadcast struct {
	RoundNb     int64      `json:"round_nb"`
	NextRoundAt *time.Time `json:"next_round_at,omitempty"`
}

// It will be sent to all clients when a round starts, with the round number and the absolute end time of the round
// (not a relative duration). Each client will calculate its own countdown using (RoundEndsAt - now), so all clients
// will converge to the same final instant regardless of their network latency at the time of reception.
type RoundStartedBroadcast struct {
	RoundNb     int64		`json:"round_nb"`
	RoundEndsAt time.Time	`json:"round_ends_at"`
	MID			int64		`json:"mid"`
	HasGuessed  bool		`json:"has_guessed"`
}

// It shares the results of the round to all players, even those who did not guess (they will see their own result as "not guessed").
// The results are sorted by player ID for consistency across clients.
// It is shared at the end of a round, with an indication of whether it was the last round or not :
// 	-> multiplayer, not the last round: NextRoundAt is provided as an absolute timestamp for when the next round will start.
// 	-> solo, not the last round: NextRoundAt is omitted, the player starts the next round himself with a "start_round" action.
// It will be nil or omitted if IsLastRound is true.
// RoundNb is the number of the round that just ended (the next one, if any, is RoundNb + 1).
type RoundResultBroadcast struct {
	RoundNb     int64                           `json:"round_nb"`
	Results     []ws_handlers.PlayerRoundResult `json:"results"`
	IsLastRound bool                            `json:"is_last_round"`
	NextRoundAt *time.Time                      `json:"next_round_at,omitempty"`
}

// encapsulation of the message sent to the client, with an action and a payload interface{}
// that can be any struct (RoundStartedBroadcast, RoundResultBroadcast, etc.)
type WsMessage struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

// roundPhase describes what h.roundState[gameId] currently represents
// so a reconnecting player gets the right kind of message (a running round vs the pause between rounds)
type roundPhase string

const (
	phaseStarting roundPhase = "starting" // start_game accepted, round 1 not started yet: in multiplayer look at NextRoundAt (countdown)
	phasePlaying  roundPhase = "playing"  // a round is currently running: look at RoundEndsAt
	phasePause    roundPhase = "pause"    // between two rounds (or after the last one): look at LastResults (+ NextRoundAt in multiplayer)
)

// to be checked when a client reconnects mid-game - refresh, dropped wifi, joining a game already in progress
// It exists from the moment the host's start_game is accepted: no state = the game hasn't started yet (still in the lobby).
type gameRoundState struct {
	Phase       roundPhase
	Solo        bool                            // decided once by start_game (only one player connected to the WS), carried over by every state change
	RoundNb     int64                           // 0 while Phase == Starting (no round played yet)
	RoundEndsAt time.Time                       // valid when Phase == Playing
	NextRoundAt time.Time                       // multiplayer only: valid when Phase == Starting, or Phase == Pause and !IsLastRound (zero otherwise)
	IsLastRound bool
	LastResults []ws_handlers.PlayerRoundResult // results of the round that just ended, valid when Phase == phasePause
	MID         int64                           // the Mapillary ID of the location for the current round, valid when Phase == Playing
}

// DataInfos carries one action through the hubw
// it can comes from a real client or be generated internally by a timer callback
type DataInfos struct {
	Action	string
	Payload json.RawMessage
	// Client is nil for events generated internally by the hub
	Client *ClientHub
}

type Hub struct {
	Clients			map[int64]map[*ClientHub]bool
	Register		chan *ClientHub
	Unregister		chan *ClientHub
	Action			chan *DataInfos
	// ONE single timer per game (map[gameId]): the round timeout while a round runs, or the automatic start of the next
	// round during a multiplayer pause. It is only a trigger: deadlines always come from rounds.started_at + game.timer,
	// and the last round from game.round_total
	GameTimers		map[int64]*time.Timer
	RoundState		map[int64]*gameRoundState // means map[gameId]currentState, used to resync reconnecting players
	LocationRepo	*repositories.LocationRepository
	LobbyRepo		*repositories.LobbyRepository
	GameRepo		*repositories.GameRepository
	UserRepo		*repositories.UserRepository
	LobbyWsHandler	*ws_handlers.LobbyWsHandler
	GuessWsHandler	*ws_handlers.GuessWsHandler
}

func NewHub(
	repo *repositories.Repositories,
	wsHandlers *ws_handlers.WsHandlers,
) *Hub {
	return &Hub{Clients:	make(map[int64]map[*ClientHub]bool),
		Register:			make(chan *ClientHub),
		Unregister:			make(chan *ClientHub),
		Action:				make(chan *DataInfos),
		GameTimers:			make(map[int64]*time.Timer),
		RoundState:			make(map[int64]*gameRoundState),
		LocationRepo:		repo.Location,
		LobbyRepo:			repo.Lobby,
		GameRepo:			repo.Game,
		UserRepo:			repo.User,
		LobbyWsHandler:		wsHandlers.Lobby,
		GuessWsHandler:		wsHandlers.Guess,
	}
}

func NewClientHub(connx *client_ws.Connexion,
    userId int64, lobbyId int64, gameId int64) *ClientHub {
    return &ClientHub{Connx: connx,
    UserId: userId,
    LobbyId: lobbyId,
    GameId: gameId}
}

// broadcast serializes the payload once and sends it to all currently connected clients for gameId
// -> serialization means that the payload is converted to JSON format before being sent to the clients
// -> deserialization means that the backend will convert the JSON back to the original data structure when it receives it
func (h *Hub) broadcast(gameId int64, payload interface{}) {
	serializedData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	for client := range h.Clients[gameId] {
		client.Connx.Send <- serializedData
	}
}

// sendToClient serializes the payload once and sends it to ONE single client 
// unlike broadcast, it is used for messages that only concern this client (errors, resync of a reconnecting player)
func (h *Hub) sendToClient(client *ClientHub, payload interface{}) {
	serializedData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	client.Connx.Send <- serializedData
}

func (h *Hub) sendToClientByUserID(gameId int64, userId int64, payload interface{}) {
    serializedData, err := json.Marshal(payload)
    if err != nil {
        log.Printf("error : %s\n", err)
        return
    }
	for client := range h.Clients[gameId] {
		if client.UserId == userId {
			client.Connx.Send <- serializedData
		}	
	}
}

// Tells ONE client why its action was refused (the other players don't need to know)
func (h *Hub) sendError(client *ClientHub, message string) {
	h.sendToClient(client, WsMessage{Action: "error", Payload: gin.H{"message": message}})
}

func (h *Hub) clientLobbyBroadcast(gameId int64, ctx context.Context) {
	lobbies, err := h.LobbyRepo.AllPlayers(ctx, gameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	usersLobby := make([]playerInfos, 0, len(lobbies))
	for _, lobby := range lobbies {
		user, err := h.UserRepo.User(ctx, lobby.UserId)
		if err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		usersLobby = append(usersLobby, playerInfos{
			ID:       user.ID,
			Username: user.Username,
			Avatar:   user.Avatar,
			Host:     lobby.IsCreator,
		})
	}
	h.broadcast(gameId, WsMessage{Action: "lobby_players", Payload: usersLobby})
}

// resyncClient sends a reconnecting player whatever the hub knows about the game
func (h *Hub) resyncClient(ctx context.Context, client *ClientHub) {
	state, ok := h.RoundState[client.GameId]
	if !ok {
		// the game hasn't started yet (still in the lobby): nothing to resync
		return
	}
	switch state.Phase {
	case phaseStarting:
		h.sendToClient(client, gameStartedMessage(state))
	case phasePlaying:
		hasGuessed, err := h.GuessWsHandler.HasGuessed(ctx, client.GameId, state.RoundNb, client.UserId)
		if err != nil {
			log.Printf("error : %s\n", err)
		}
		h.sendToClient(client, WsMessage{
			Action: "round_started",
			Payload: RoundStartedBroadcast{
				RoundNb:		state.RoundNb,
				RoundEndsAt:	state.RoundEndsAt.UTC(),
				MID:			state.MID,
				HasGuessed:		hasGuessed,
			},
		})
	case phasePause:
		payload := RoundResultBroadcast{RoundNb: state.RoundNb, Results: state.LastResults, IsLastRound: state.IsLastRound}
		// countdown to the next round: multiplayer only 
		if !state.IsLastRound && !state.Solo {
			nextRoundAt := state.NextRoundAt
			payload.NextRoundAt = &nextRoundAt
		}
		h.sendToClient(client, WsMessage{Action: "round_results", Payload: payload})
	}
}

// inProgress tells whether gameId has started and isn't over yet (its last round hasn't ended)
func (h *Hub) inProgress(gameId int64) bool {
	state, ok := h.RoundState[gameId]
	return ok && !(state.Phase == phasePause && state.IsLastRound)
}

// Checks whether the last round of gameId has ended
func (h *Hub) gameOver(gameId int64) bool {
	state, ok := h.RoundState[gameId]
	return ok && state.Phase == phasePause && state.IsLastRound
}

// It decides, from the hub's state only (never from anything a client sent), which round the game is waiting to start:
//   - phaseStarting: start_game has been accepted, so it's round 1
//   - phasePause: we're between two rounds, so it's the next one (none after the last round)
//   - no state (the game hasn't started: use start_game) or phasePlaying (a round is already running): nothing to start
func (h *Hub) nextRoundToStart(gameId int64) (int64, bool) {
	state, ok := h.RoundState[gameId]
	if !ok {
		return 0, false
	}
	if state.Phase == phaseStarting {
		return 1, true
	}
	if state.Phase == phasePause && !state.IsLastRound {
		return state.RoundNb + 1, true
	}
	return 0, false
}

// Checks whether gameId is played in solo
func (h *Hub) isSolo(gameId int64) bool {
	state, ok := h.RoundState[gameId]
	return ok && state.Solo
}

// Checks whether a game has started already
func (h *Hub) isStarted(ctx context.Context, gameId int64) bool {
	if _, ok := h.RoundState[gameId]; ok {
		return true
	}
	_, err := h.GuessWsHandler.ResolveRoundID(ctx, gameId, 1)
	if errors.Is(err, repositories.ErrRoundNotFound) {
		return false
	}
	if err != nil {
		log.Printf("error : %s\n", err)
	}
	// round 1 exists, or the database couldn't tell: refuse rather than risk starting the game twice
	return true
}

func (h *Hub) registerClient(client *ClientHub) {
	if _, ok := h.Clients[client.GameId]; !ok {
		h.Clients[client.GameId] = make(map[*ClientHub]bool)
	}
	h.Clients[client.GameId][client] = true
}

func (h *Hub) unregisterClient(client *ClientHub) {
	if _, ok := h.Clients[client.GameId]; ok {
		delete(h.Clients[client.GameId], client)
	}
	if len(h.Clients[client.GameId]) == 0 {
		delete(h.Clients, client.GameId)
	}
	close(client.Connx.Send)
}

// Routes an action received (from a client or generated by the hub) to the appropriate handler.
// Handled sequentially (only one goroutine reads h.action in Run) : no race condition possible between
// two actions for the same gameId, even if they come from different clients.
func (h *Hub) handleAction(ctx context.Context, data *DataInfos) {
	switch data.Action {
	case "leave":
		if h.inProgress(data.Client.GameId) {
			// during a game, a lost connection keeps the player in it (he can reconnect):
			// leaving on purpose must remove him right away
			if !h.LobbyWsHandler.LeaveLobby(ctx, data.Client.UserId) {
				log.Println("error : the client can't leave the lobby in the database")
			}
			h.clientLobbyBroadcast(data.Client.GameId, ctx)
		}
		// close every connection of the kicked player (he may have several tabs open)
		for client := range h.Clients[data.Client.GameId] {
			if client.UserId == data.Client.UserId {
				client.Connx.Conn.Close()
			}
		}

	case "kick_player":
		var targetClient kickPlayer
		if err := json.Unmarshal(data.Payload, &targetClient); err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		if !h.LobbyWsHandler.RemovePlayer(ctx, data.Client.UserId, targetClient.TargetId) {
			return
		}
		h.clientLobbyBroadcast(data.Client.GameId, ctx)
		// close every connection of the kicked player (he may have several tabs open)
		for client := range h.Clients[data.Client.GameId] {
			if client.UserId == targetClient.TargetId {
				client.Connx.Conn.Close()
			}
		}

	case "guess":
		h.handleGuess(ctx, data)

	case "start_game":
		h.handleStartGame(ctx, data.Client)

	case "start_round":
		if data.Client != nil {
			h.handleStartRound(ctx, data.Client)
			return
		}
		// client == nil > automatic start of the next round at the end of a multiplayer pause (server-generated)
		var evt internalRoundEvent
		if err := json.Unmarshal(data.Payload, &evt); err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		// only start the round the game is really waiting for (ignores a stale timer)
		if roundNb, ok := h.nextRoundToStart(evt.GameId); ok && roundNb == evt.RoundNb {
			h.StartRound(ctx, evt.GameId, evt.RoundNb)
		}

	case "round_timeout":
		if data.Client != nil {
			// round_timeout must only ever be generated internally by armTimer
			log.Printf("warning : client %d tried to trigger round_timeout directly, ignored\n", data.Client.UserId)
			return
		}
		var evt internalRoundEvent
		if err := json.Unmarshal(data.Payload, &evt); err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		h.EndRound(ctx, evt.GameId, evt.RoundNb)

	case "finished":
		// each player leaves the finished game on his own
		// (the others stay on the results as long as they want)
		// and only once the last round has ended and its results have been sent
		if data.Client == nil {
			return
		}
		if !h.gameOver(data.Client.GameId) {
			h.sendError(data.Client, "the game is not finished yet")
			return
		}
		if !h.LobbyWsHandler.LeaveFinishedGame(ctx, data.Client.UserId, data.Client.GameId) {
			h.sendError(data.Client, "could not leave the game, try again")
			return
		}
		// close every connection of this player only (he may have several tabs open)
		for client := range h.Clients[data.Client.GameId] {
			if client.UserId == data.Client.UserId {
				client.Connx.Conn.Close()
			}
		}
	default:
		log.Printf("unknown action : %s\n", data.Action)
	}
}

// handleGuess persists the received guess, then ends the round if all active players in the game have guessed
func (h *Hub) handleGuess(ctx context.Context, data *DataInfos) {
	// a guess only counts while a round is running. During a pause the results (so the right location) are on screen,
	// and a late guess would still earn the XP of a perfect guess (base_score doesn't depend on the time left)
	if state, ok := h.RoundState[data.Client.GameId]; !ok || state.Phase != phasePlaying {
		h.sendError(data.Client, "no round is running")
		return
	}
	createdGuess, ok := h.GuessWsHandler.AddGuess(ctx, data.Payload, data.Client.UserId, data.Client.GameId)
	if !ok {
		// TODO : renvoyer une erreur ciblée uniquement à data.Client
		return
	}

	activePlayers, err := h.LobbyRepo.AllPlayers(ctx, data.Client.GameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}

	roundGuesses, err := h.GuessWsHandler.GetRoundGuesses(ctx, createdGuess.RoundID)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}

	if len(roundGuesses) >= len(activePlayers) {
		h.EndRound(ctx, data.Client.GameId, createdGuess.RoundNb)
	}
}

// Returns the distinct players currently connected to the WS for gameId
// (a player with several tabs open, so several connections, appears once)
func (h *Hub) connectedPlayers(gameId int64) map[int64]bool {
	users := make(map[int64]bool)
	for client := range h.Clients[gameId] {
		users[client.UserId] = true
	}
	return users
}

// handleStartGame starts the game, at the host's request and only once. The mode is decided here, once, from the number
// of players connected to the game's WS, and never changes afterwards (even if someone joins or leaves during the game):
//   - solo (only one player connected): round 1 starts right away, then each next round is started by the player ("start_round")
//   - multiplayer: every round, the first one included, starts automatically after a pause of preRoundPause
// The lobby closes at the same time: nobody can join the game afterwards (its players can still reconnect).
func (h *Hub) handleStartGame(ctx context.Context, client *ClientHub) {
	gameId := client.GameId
	if !h.LobbyWsHandler.IsHost(ctx, client.UserId, gameId) {
		h.sendError(client, "only the host can start the game")
		return
	}
	if h.isStarted(ctx, gameId) {
		h.sendError(client, "the game has already started")
		return
	}
	connected := h.connectedPlayers(gameId)
	// the lobby closes: its join code is removed, and whoever is in the lobby without being connected is removed from it
	if !h.LobbyWsHandler.CloseLobby(ctx, gameId, connected) {
		h.sendError(client, "the game could not be started, try again")
		return
	}
	h.clientLobbyBroadcast(gameId, ctx)

	// from here on the game counts as started (see isStarted): any other start_game will be refused
	state := &gameRoundState{Phase: phaseStarting, Solo: len(connected) == 1}
	if !state.Solo {
		state.NextRoundAt = time.Now().Add(preRoundPause)
	}
	h.RoundState[gameId] = state
	h.broadcast(gameId, gameStartedMessage(state))

	if state.Solo {
		h.StartRound(ctx, gameId, 1)
	} else {
		h.armTimer(gameId, time.Until(state.NextRoundAt), "start_round", 1)
	}
}

// Handles a "start_round" sent by a client: only accepted in solo, from the host,
// when the game is waiting for its next round (results of the previous one on screen)
// (In multiplayer, rounds start automatically after a pause)
func (h *Hub) handleStartRound(ctx context.Context, client *ClientHub) {
	state, ok := h.RoundState[client.GameId]
	if !ok {
		h.sendError(client, "the game has not started yet")
		return
	}
	if !state.Solo {
		h.sendError(client, "rounds start automatically in multiplayer")
		return
	}
	if !h.LobbyWsHandler.IsHost(ctx, client.UserId, client.GameId) {
		h.sendError(client, "only the host can start a round")
		return
	}
	roundNb, ok := h.nextRoundToStart(client.GameId)
	if !ok {
		// a round is already running, or the last one is over
		h.sendError(client, "there is no round to start")
		return
	}
	h.StartRound(ctx, client.GameId, roundNb)
}

// gameStartedMessage builds the "game_started" message from the game's state
func gameStartedMessage(state *gameRoundState) WsMessage {
	data := GameStartedBroadcast{RoundNb: 1}
	if !state.NextRoundAt.IsZero() {
		nextRoundAt := state.NextRoundAt.UTC()
		data.NextRoundAt = &nextRoundAt
	}
	return WsMessage{Action: "game_started", Payload: data}
}

// It arms the round timer for the given gameId and roundNb, and broadcasts the round_started event to all clients.
// The timer is set to the remaining duration of the round, calculated from the round's StartedAt timestamp and the round duration.
// If the round has already expired (remaining <= 0), it immediately ends the round instead of arming a negative timer:
// a round that already existed (ex: a start that failed halfway and was retried later) never gets more time than expected,
// and all clients stay synchronized with the actual round start time.
func (h *Hub) StartRound(ctx context.Context, gameId, roundNb int64) {
	round, err := h.GuessWsHandler.EnsureRound(ctx, gameId, roundNb)
	if err != nil {
		log.Printf("error starting round: %s\n", err)
		return
	}

	duration, err := h.GuessWsHandler.GetRoundDuration(ctx, gameId)
	if err != nil {
		log.Printf("error getting round duration: %s\n", err)
		return
	}

	MID, err := h.LocationRepo.GetMIDByID(ctx, round.LocationId)
	if err != nil {
		log.Printf("error getting MID: %s\n", err)
		return
	}

	roundEndsAt := round.StartedAt.Add(duration)
	h.RoundState[gameId] = &gameRoundState{
		Phase:       phasePlaying,
		Solo:        h.isSolo(gameId), // the mode never changes after start_game
		RoundNb:     roundNb,
		RoundEndsAt: roundEndsAt,
		MID:         MID, //TODO should be kept server side if we fetch through the proxy, not sent to clients
	}

	remaining := time.Until(roundEndsAt)
	if remaining <= 0 {
		// round has expired: end it immediately (the state is set just above because EndRound only ends the running round)
		h.EndRound(ctx, gameId, roundNb)
		return
	}

	h.armTimer(gameId, remaining, "round_timeout", roundNb)

	h.broadcast(gameId, WsMessage{
		Action: "round_started",
		Payload: RoundStartedBroadcast{
			RoundNb:		roundNb,
			RoundEndsAt:	roundEndsAt.UTC(),
			MID:			MID,
		},
	})
}

// armTimer (re)arms THE pending timer of a game: the round timeout while a round is running ("round_timeout"), or the
// automatic start of the next round during a multiplayer pause ("start_round"). There's never more than one at a time,
// so arming one stops the previous one. When it fires, the callback doesn't touch the hub's state itself: it sends an
// internal event (Client nil) on h.Action, handled sequentially in Run like any other action.
func (h *Hub) armTimer(gameId int64, delay time.Duration, action string, roundNb int64) {
	if existing, ok := h.GameTimers[gameId]; ok {
		existing.Stop()
	}
	h.GameTimers[gameId] = time.AfterFunc(delay, func() {
		payload, err := json.Marshal(internalRoundEvent{GameId: gameId, RoundNb: roundNb})
		if err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		h.Action <- &DataInfos{Action: action, Payload: payload}
	})
}

// Finalises the round roundNb: broadcasts the results of ALL active players (whether they guessed or not), then,
// unless it was the last round:
//   - multiplayer: schedules the automatic start of the next round after a pause of preRoundPause (10sec)
//   - solo: schedules nothing, the game waits for the player's "start_round" (see handleStartRound)
func (h *Hub) EndRound(ctx context.Context, gameId, roundNb int64) {
	state, ok := h.RoundState[gameId]
	if !ok || state.Phase != phasePlaying || state.RoundNb != roundNb {
		// only the round currently running can be ended. Happens when the last guess of a round is processed at almost
		// the same instant the round's timer fires: the second event arrives when the round is already over
		return
	}
	solo := state.Solo
	if timer, ok := h.GameTimers[gameId]; ok {
		timer.Stop()
		delete(h.GameTimers, gameId)
	}

	roundID, err := h.GuessWsHandler.ResolveRoundID(ctx, gameId, roundNb)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	activePlayers, err := h.LobbyRepo.AllPlayers(ctx, gameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	roundGuesses, err := h.GuessWsHandler.GetRoundGuesses(ctx, roundID)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	results, err := h.GuessWsHandler.BuildRoundResults(ctx, gameId, roundID, roundNb, activePlayers, roundGuesses)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	isLastRound, err := h.GuessWsHandler.IsLastRound(ctx, gameId, roundNb)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	if isLastRound {
		// the game is over: finished_at is set and every player gets the base_score of this game in his xp (only once)
		if err := h.GameRepo.SetFinished(ctx, gameId); err != nil {
			log.Printf("error : %s\n", err)
		}
		// the players who are still disconnected leave the game now, otherwise their lobby would stay open and they
		// couldn't create or join another one (the connected players leave when they close the game, see Unregister)
		h.LobbyWsHandler.CloseLobby(ctx, gameId, h.connectedPlayers(gameId))
	}

	broadcastPayload := RoundResultBroadcast{RoundNb: roundNb, Results: results, IsLastRound: isLastRound}
	var nextRoundAt time.Time
	if !isLastRound && !solo {
		// calculated once: this value is the reference. The timer armed below is based on the remaining time
		// until this same instant (time.Until), so the next round starts exactly when all clients have seen their countdown reach zero.
		nextRoundAt = time.Now().Add(preRoundPause)
		broadcastPayload.NextRoundAt = &nextRoundAt
	}
	h.RoundState[gameId] = &gameRoundState{
		Phase:       phasePause,
		Solo:        solo,
		RoundNb:     roundNb,
		NextRoundAt: nextRoundAt,
		IsLastRound: isLastRound,
		LastResults: results,
	}

	h.broadcast(gameId, WsMessage{Action: "round_results", Payload: broadcastPayload})

	if isLastRound && len(h.Clients[gameId]) == 0 {
		// the game ended with nobody connected (everyone lost his connection): no disconnection will come
		// to drop its state later (see Run), so it goes now
		delete(h.RoundState, gameId)
	}
	if isLastRound || solo {
		// solo: no timer, the next round will be started by the player's "start_round" (see handleStartRound)
		return
	}
	h.armTimer(gameId, time.Until(nextRoundAt), "start_round", roundNb+1)
}

func JoinLobby(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "The userId is empty"})
			return
		}        
		// Check if the player is in the right lobby
        lobby, err := hub.LobbyRepo.ActualLobby(c.Request.Context(), userId.(int64))
		if err != nil {
            log.Printf("error : %s\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"message": "could not verify lobby"})
            return
        }
        if lobby == nil {
            // Previously this returned with no response at all: the client's upgrade request was left hanging
            // so this message avoids a hung connection
            c.JSON(http.StatusForbidden, gin.H{"message": "you are not in a lobby"})
            return
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("error : %s\n", err)
			return 
		}
		connx := client_ws.NewConnexion(conn)
		clientHub := NewClientHub(connx, userId.(int64),
			lobby.ID, lobby.GameId)
		hub.Register <- clientHub
		go clientHub.Connx.WritePump()
		actionClient := func(data []byte) {
			var dataInfos DataInfos
			err := json.Unmarshal(data, &dataInfos)
			if err != nil {
				log.Printf("error : %s\n", err)
				return
			}
			dataInfos.Client = clientHub
			hub.Action <- &dataInfos 
		}
		unregisterClient := func() {
			hub.Unregister <- clientHub
		}
		clientHub.Connx.ReadPump(actionClient, unregisterClient)
	}
}

// handleDisconnection decides what a closed connection means for the player, depending on where the game is
func (h *Hub) handleDisconnection(ctx context.Context, client *ClientHub) {
	if h.connectedPlayers(client.GameId)[client.UserId] {
		// the player still has another tab open on this game: he hasn't left
		return
	}
	if h.inProgress(client.GameId) {
		// connection lost during the game: the player stays in it and can reconnect
		// He only leaves with a "leave" action, or at the end of the game if he hasn't come back
		return
	}
	if h.gameOver(client.GameId) {
		// after the last round the player just leaves: no host or join code to hand over
		h.LobbyWsHandler.LeaveFinishedGame(ctx, client.UserId, client.GameId)
		return
	}
	// still in the lobby: the player leaves it (the host and the join code go to another player if needed)
	if !h.LobbyWsHandler.LeaveLobby(ctx, client.UserId) {
		log.Println("error : the client can't leave the lobby in the database")
		return
	}
	h.clientLobbyBroadcast(client.GameId, ctx)
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <- h.Register:
			// only the game's players get here (JoinLobby requires an open lobby row, and nobody can join a started game):
			// during a game, it's a player coming back after losing his connection
			h.registerClient(client)
			h.clientLobbyBroadcast(client.GameId, ctx)
			if data, ok := h.LobbyWsHandler.SendInfos(ctx, client.GameId); ok {
				h.sendToClient(client, WsMessage{Action: "lobby_infos", Payload: data})
			}
			// then where the game is (countdown, running round or results), so a reconnecting player is back in his round
			h.resyncClient(ctx, client)
		case client := <-h.Unregister:
			h.unregisterClient(client)
			h.handleDisconnection(ctx, client)
			// a finished game keeps its state while players are still on its results (gameOver reads it),
			// and loses it when the last one is gone: nobody can come back to a finished game
			if len(h.Clients[client.GameId]) == 0 && h.gameOver(client.GameId) {
				delete(h.RoundState, client.GameId)
			}
		case serializedData := <- h.Action:
			h.handleAction(ctx, serializedData)
		}
	}
}