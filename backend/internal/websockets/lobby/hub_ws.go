package lobby_ws

import (
	"os"
	"backend/internal/repositories"
	"backend/internal/tools/mapillary"
	securesend "backend/internal/tools/secure_send"
	"backend/internal/websockets/client"
	"backend/internal/websockets/handlers"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
		if r.Header.Get("Origin") != os.Getenv("PUBLIC_URL") {
			return false
		}
		return true
	},
}

// multiplayer only: pause before EVERY round
const preRoundPause = 10 * time.Second

// when the database fails at the end of a round, EndRound tries again every endRoundRetryDelay
// at most endRoundMaxRetries times (about a minute), after that the game is stopped
const endRoundRetryDelay = 2 * time.Second
const endRoundMaxRetries = 30

// the image of a round that players can't load is checked on Mapillary at most maxLocationChecks times per round:
// a front that keeps failing on its side (connection, browser) would otherwise ask again after each reload
const maxLocationChecks = 3

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

// Carries the game_id/round_nb for events the hub itself generates
type internalRoundEvent struct {
	GameId  int64 `json:"game_id"`
	RoundNb int64 `json:"round_nb"`
}

// Sent by a player whose front can't load the image of the running round (black screen)
type locationErrorPayload struct {
	RoundNb int64 `json:"round_nb"`
}

// Sent back to the hub by the goroutine that asked Mapillary about the image of a round (see handleLocationError)
type locationCheckedEvent struct {
	GameId	int64	`json:"game_id"`
	RoundNb	int64	`json:"round_nb"`
	MID		int64	`json:"mid"`
	Broken	bool	`json:"broken"`
}

// What the hub needs from Mapillary (see tools/mapillary)
type imageChecker interface {
	ImageExists(ctx context.Context, mid int64) (bool, error)
}

// Sent to all clients when the host starts the game
// Resent to a client that reconnects during the countdown before round 1
type GameStartedBroadcast struct {
	RoundNb     int64      `json:"round_nb"`
	NextRoundAt *time.Time `json:"next_round_at,omitempty"`
}

// Sent to all clients when a round starts
type RoundStartedBroadcast struct {
	RoundNb     int64		`json:"round_nb"`
	RoundEndsAt time.Time	`json:"round_ends_at"`
	MID			int64		`json:"mid"`
	HasGuessed  bool		`json:"has_guessed"`
}

// Shares the results of EVERY round played so far, in decreasing order
type GameResultsBroadcast struct {
	Rounds		[]ws_handlers.RoundResult	`json:"rounds"`
	NextRoundAt	*time.Time					`json:"next_round_at,omitempty"`
}

// encapsulation of the message sent to the client, with an action and a payload interface{}
// that can be any struct (GameResultsBroadcast, etc.)
type WsMessage struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

// Describes what h.roundState[gameId] currently represents
// so a reconnecting player gets the right kind of message (a running round vs the pause between rounds)
type roundPhase string

const (
	phasePlaying roundPhase = "playing"
	phasePause   roundPhase = "pause"
)

// to be checked when a client reconnects mid-game - refresh, dropped wifi, joining a game already in progress
// It exists from the moment the host's start_game is accepted: no state = the game hasn't started yet or is finished
type gameRoundState struct {
	Phase		roundPhase // phasePlaying, phasePause, or ""
	Solo		bool // decided once by start_game (only one player connected to the WS), carried over by every state change
	RoundNb		int64 // round running (Playing) or last round ended (Pause, game over); 0 in the pause before round 1
	RoundEndsAt	time.Time // valid when Phase == Playing
	EndRetries	int // failed attempts to end the running round so far
	NextRoundAt	time.Time // multiplayer only, valid when Phase == Pause: when the next round starts (zero otherwise)
	GameResults	[]ws_handlers.RoundResult // results of every round played so far, newest first
	MID			int64 // the Mapillary ID of the location for the current round, valid when Phase == Playing
	CheckingLocation	bool // Mapillary is being asked whether MID still exists (see handleLocationError)
	WaitingReload		map[*ClientHub]bool // clients that couldn't load MID, waiting for the answer of Mapillary
	LocationChecks		int // Mapillary checks already done for this round (at most maxLocationChecks)
}

// Carries one action through the hubw
// it can comes from a real client or be generated internally by a timer callback
// Client is nil for events generated internally by the hub
type DataInfos struct {
	Action	string
	Payload json.RawMessage
	Client *ClientHub
}

type Hub struct {
	Clients			map[int64]map[*ClientHub]bool
	Register		chan *ClientHub
	Unregister		chan *ClientHub
	Action			chan *DataInfos
	GameTimers		map[int64]*time.Timer
	RoundState		map[int64]*gameRoundState // means map[gameId]currentState, used to resync reconnecting players
	LocationRepo	*repositories.LocationRepository
	LobbyRepo		*repositories.LobbyRepository
	GameRepo		*repositories.GameRepository
	UserRepo		*repositories.UserRepository
	StatsRepo		*repositories.StatisticsRepository
	LobbyWsHandler	*ws_handlers.LobbyWsHandler
	GuessWsHandler	*ws_handlers.GuessWsHandler
	Mapillary		imageChecker // nil without MAPILLARY_ACCESS_TOKEN: an image that can't be loaded is then only reloaded
}

func NewHub(
	repo *repositories.Repositories,
	wsHandlers *ws_handlers.WsHandlers,
) *Hub {
	hub := &Hub{Clients:	make(map[int64]map[*ClientHub]bool),
		Register:			make(chan *ClientHub),
		Unregister:			make(chan *ClientHub),
		Action:				make(chan *DataInfos),
		GameTimers:			make(map[int64]*time.Timer),
		RoundState:			make(map[int64]*gameRoundState),
		LocationRepo:		repo.Location,
		LobbyRepo:			repo.Lobby,
		GameRepo:			repo.Game,
		UserRepo:			repo.User,
		StatsRepo:			repo.Statistics,
		LobbyWsHandler:		wsHandlers.Lobby,
		GuessWsHandler:		wsHandlers.Guess,
	}
	// only set when there is a token: a nil *mapillary.Checker stored in the interface wouldn't be nil anymore
	checker, err := mapillary.NewChecker()
	if err != nil {
		log.Printf("warning : %s, an image that a player can't load will be reloaded, never replaced\n", err)
		return hub
	}
	hub.Mapillary = checker
	return hub
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
		select {
		case client.Connx.Send <- serializedData:
		default:
			log.Printf("client_error : the client %d can't read\n", client.UserId)
		}
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
	select {
	case client.Connx.Send <- serializedData:
	default:
		log.Printf("client_error : the client %d can't read\n", client.UserId)
	}
}

func (h *Hub) sendToClientByUserID(gameId int64, userId int64, payload interface{}) {
    serializedData, err := json.Marshal(payload)
    if err != nil {
        log.Printf("error : %s\n", err)
        return
    }
	for client := range h.Clients[gameId] {
		if client.UserId == userId {
			select {
			case client.Connx.Send <- serializedData:
			default:
				log.Printf("client_error : the client %d can't read\n", client.UserId)
			}
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
		if state.RoundNb == 0 {
			// pause before round 1: the countdown of the start of the game (multiplayer)
			h.sendToClient(client, gameStartedMessage(state))
			return
		}
		if len(state.GameResults) == 0 {
			return
		}
		payload := GameResultsBroadcast{Rounds: state.GameResults}
		// countdown to the next round: multiplayer only (in solo the player starts it himself with start_round)
		if !state.Solo {
			nextRoundAt := state.NextRoundAt
			payload.NextRoundAt = &nextRoundAt
		}
		h.sendToClient(client, WsMessage{Action: "round_results", Payload: payload})
	default:
		// no phase: the game is over (for instance a second tab of a player still on the final results)
		if len(state.GameResults) == 0 {
			// stopped by abortGame: no final results, the player is sent out like the others
			closeWithServerError(client)
			return
		}
		h.sendToClient(client, WsMessage{
			Action:  "game_results",
			Payload: GameResultsBroadcast{Rounds: state.GameResults},
		})
	}
}

// inProgress tells whether gameId has started and isn't over yet (its last round hasn't ended)
func (h *Hub) inProgress(gameId int64) bool {
	state, ok := h.RoundState[gameId]
	return ok && state.Phase != ""
}

// Checks whether the last round of gameId has ended
func (h *Hub) gameOver(gameId int64) bool {
	state, ok := h.RoundState[gameId]
	return ok && state.Phase == ""
}

// It decides, from the hub's state only (never from anything a client sent), which round the game is waiting to start:
// -phasePause: it's the next one (round 1 in the pause before the first round)
// -ELSE: during a round (phasePlaying) or once the game is over (no phase), there is nothing to start
func (h *Hub) nextRoundToStart(gameId int64) (int64, bool) {
	state, ok := h.RoundState[gameId]
	if !ok || state.Phase != phasePause {
		return 0, false
	}
	return state.RoundNb + 1, true
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
	// round 1 exists, so we do not start the game twice
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
			// during a game, a lost connection keeps the player in it (reconnections allowed):
			// leaving on purpose removes the player right away (reconnections not allowed)
			if !h.LobbyWsHandler.LeaveLobby(ctx, data.Client.UserId) {
				log.Println("error : the client can't leave the lobby in the database")
			}
			h.clientLobbyBroadcast(data.Client.GameId, ctx)
			h.dropIfAbandoned(ctx, data.Client.GameId)
		}
		// close every connection of the player (all tabs opened)
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
		// close every connection of the kicked player (all tabs open)
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
			// round_timeout can only ever be generated internally by armTimer
			log.Printf("warning : client %d tried to trigger round_timeout directly, ignored\n", data.Client.UserId)
			return
		}
		var evt internalRoundEvent
		if err := json.Unmarshal(data.Payload, &evt); err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		// when a round starts again with a new location (see handleLocationChecked), the timer of the old one can fire
		// right before being stopped: its event arrives before the new end of the round. A real timeout never comes
		// early (and the retries of EndRound are counted in EndRetries), so an early one only re-arms the timer
		if state, ok := h.RoundState[evt.GameId]; ok && state.Phase == phasePlaying && state.RoundNb == evt.RoundNb &&
			state.EndRetries == 0 && time.Until(state.RoundEndsAt) > 0 {
			h.armTimer(ctx, evt.GameId, time.Until(state.RoundEndsAt), "round_timeout", evt.RoundNb)
			return
		}
		h.EndRound(ctx, evt.GameId, evt.RoundNb)

	case "location_error":
		// the front of this player can't load the image of the running round
		if data.Client == nil {
			return
		}
		h.handleLocationError(ctx, data)

	case "location_checked":
		if data.Client != nil {
			// location_checked can only ever be generated internally (see handleLocationError):
			// a client could otherwise have any location replaced
			log.Printf("warning : client %d tried to trigger location_checked directly, ignored\n", data.Client.UserId)
			return
		}
		h.handleLocationChecked(ctx, data.Payload)

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
		// close every connection of this player only (all tabs open)
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
	// a guess only counts while a round is running
	if state, ok := h.RoundState[data.Client.GameId]; !ok || state.Phase != phasePlaying {
		h.sendError(data.Client, "no round is running")
		return
	}
	createdGuess, ok := h.GuessWsHandler.AddGuess(ctx, data.Payload, data.Client.UserId, data.Client.GameId)
	if !ok {
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

	// only the guesses of the players still in the game count: a player who guessed then left (or was kicked)
	// must not stand in for one who hasn't guessed yet
	inGame := make(map[int64]bool, len(activePlayers))
	for _, player := range activePlayers {
		inGame[player.UserId] = true
	}
	guessed := 0
	for _, g := range roundGuesses {
		if inGame[g.UserID] {
			guessed++
		}
	}
	if guessed >= len(activePlayers) {
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
	// pause before round 1 (RoundNb 0): 10s of countdown in multiplayer, round 1 starts right away in solo
	state := &gameRoundState{Phase: phasePause, Solo: len(connected) == 1}
	if !state.Solo {
		state.NextRoundAt = time.Now().Add(preRoundPause)
	}
	h.RoundState[gameId] = state
	h.broadcast(gameId, gameStartedMessage(state))

	if state.Solo {
		h.StartRound(ctx, gameId, 1)
	} else {
		h.armTimer(ctx, gameId, time.Until(state.NextRoundAt), "start_round", 1)
	}
}

// Handles a "start_round" sent by a client: only accepted in solo, from the host,
// when the game is waiting for its next round (results of the previous one on screen)
// "start_round" is not permitted in multiplayer
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
		MID:         MID,
	}
	remaining := time.Until(roundEndsAt)
	if remaining <= 0 {
		// round has expired: end it immediately (the state is set just above because EndRound only ends the running round)
		h.EndRound(ctx, gameId, roundNb)
		return
	}

	h.armTimer(ctx, gameId, remaining, "round_timeout", roundNb)

	h.broadcast(gameId, WsMessage{
		Action: "round_started",
		Payload: RoundStartedBroadcast{
			RoundNb:		roundNb,
			RoundEndsAt:	roundEndsAt.UTC(),
			MID:			MID,
		},
	})
}

// A player can't load the image of the running round (black screen). Mapillary is asked whether the image still exists,
// once at a time however many players report it (the others wait for the same answer), at most maxLocationChecks times
// per round. The call can take a few seconds: it never runs in the hub's goroutine, which would freeze every game
// meanwhile. Like a timer, the goroutine doesn't touch the hub's state: it sends the answer back as an internal event
func (h *Hub) handleLocationError(ctx context.Context, data *DataInfos) {
	var report locationErrorPayload
	if err := json.Unmarshal(data.Payload, &report); err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	state, ok := h.RoundState[data.Client.GameId]
	// only for the round running now: a report sent just before the end of the round is ignored
	if !ok || state.Phase != phasePlaying || state.RoundNb != report.RoundNb {
		return
	}
	if !state.CheckingLocation && state.LocationChecks >= maxLocationChecks {
		log.Printf("warning : the image of game %d round %d was already checked %d times, report of client %d ignored\n",
			data.Client.GameId, state.RoundNb, state.LocationChecks, data.Client.UserId)
		return
	}
	if state.WaitingReload == nil {
		state.WaitingReload = make(map[*ClientHub]bool)
	}
	state.WaitingReload[data.Client] = true
	if state.CheckingLocation {
		// already being checked: this client gets the same answer as the others
		return
	}
	state.CheckingLocation = true
	state.LocationChecks++

	evt := locationCheckedEvent{GameId: data.Client.GameId, RoundNb: state.RoundNb, MID: state.MID}
	checker := h.Mapillary
	go func() {
		defer func() {
			// out of the hub's goroutine, safely doesn't protect this one: a panic here would stop the whole server
			if r := recover(); r != nil {
				log.Printf("panic in the location check of game %d: %v\n%s", evt.GameId, r, debug.Stack())
			}
		}()
		if checker != nil {
			exists, err := checker.ImageExists(ctx, evt.MID)
			if err != nil {
				log.Printf("error : could not check the mapillary image %d (game %d round %d): %s\n",
					evt.MID, evt.GameId, evt.RoundNb, err)
			}
			evt.Broken = err == nil && !exists
		}
		payload, err := json.Marshal(evt)
		if err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		securesend.SecureSend(ctx, h.Action, &DataInfos{Action: "location_checked", Payload: payload},
			fmt.Sprintf("Hub : the location check of game %d can't be delivered", evt.GameId))
	}()
}

// Answer of Mapillary about the image of a round (see handleLocationError):
//   - the image exists (or Mapillary couldn't tell): the problem is on the side of the players who reported it,
//     only them reload the round, with the same location
//   - the image doesn't exist anymore: the location is disabled and replaced by a new one, and the round starts
//     again from the beginning for every player
func (h *Hub) handleLocationChecked(ctx context.Context, payload json.RawMessage) {
	var evt locationCheckedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	state, ok := h.RoundState[evt.GameId]
	// the round ended, or its location already changed, while Mapillary was answering: nothing left to do
	if !ok || state.Phase != phasePlaying || state.RoundNb != evt.RoundNb || state.MID != evt.MID {
		return
	}
	waiting := state.WaitingReload
	state.CheckingLocation = false
	state.WaitingReload = nil

	if !evt.Broken {
		log.Printf("the image of game %d round %d (mapillary %d) can't be loaded by %d client(s) but isn't known as broken: they reload it\n",
			evt.GameId, evt.RoundNb, evt.MID, len(waiting))
		for client := range waiting {
			// a client disconnected meanwhile is unregistered: its channel is closed, nothing can be sent to it anymore
			if _, ok := h.Clients[evt.GameId][client]; ok {
				h.resyncClient(ctx, client)
			}
		}
		return
	}

	oldLocationId, newLocationId, err := h.GameRepo.ReplaceRoundLocation(ctx, evt.GameId, evt.RoundNb)
	if err != nil {
		// the players keep their black screen until the end of the round, the game goes on with the next one
		log.Printf("error : the location of game %d round %d (mapillary %d) is broken but could not be replaced: %s\n",
			evt.GameId, evt.RoundNb, evt.MID, err)
		return
	}
	log.Printf("location %d (mapillary %d) can't be loaded anymore: disabled, and replaced by location %d in game %d round %d\n",
		oldLocationId, evt.MID, newLocationId, evt.GameId, evt.RoundNb)
	// the round starts again: new end, new timer, new MID in the state, and round_started to every player
	h.StartRound(ctx, evt.GameId, evt.RoundNb)
}

// armTimer (re)arms the pending timer of a game: the round timeout while a round is running ("round_timeout")
// or the automatic start of the next round during a multiplayer pause ("start_round")
// There's never more than one at a time, so arming one stops the previous one.
// When it fires, the callback doesn't touch the hub's state itself: it sends an internal event (Client nil) on h.Action
func (h *Hub) armTimer(ctx context.Context, gameId int64, delay time.Duration,
	action string, roundNb int64) {
	if existing, ok := h.GameTimers[gameId]; ok {
		existing.Stop()
	}
	h.GameTimers[gameId] = time.AfterFunc(delay, func() {
		payload, err := json.Marshal(internalRoundEvent{GameId: gameId, RoundNb: roundNb})
		if err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		securesend.SecureSend(ctx, h.Action, &DataInfos{Action: action, Payload: payload},
			fmt.Sprintf("Hub : the %s event of game %d can't be delivered", action, gameId))
	})
}

// This closes one connection with the websocket close code 1011, the "internal server error" of websockets
func closeWithServerError(client *ClientHub) {
	message := websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "internal server error")
	if err := client.Connx.Conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second)); err != nil {
		log.Printf("error : %s\n", err)
	}
	client.Connx.Conn.Close()
}

// It stops a game that can't go on (EndRound kept failing)
func (h *Hub) abortGame(ctx context.Context, gameId int64, roundNb int64, solo bool) {
	// no phase and no results: the game is over and was stopped
	h.RoundState[gameId] = &gameRoundState{Solo: solo, RoundNb: roundNb}
	// every player leaves the game, connected or not
	// (if the database answers again: otherwise each one leaves when his closed connection is unregistered)
	h.LobbyWsHandler.CloseLobby(ctx, gameId, nil)
	for client := range h.Clients[gameId] {
		closeWithServerError(client)
	}
	if len(h.Clients[gameId]) == 0 {
		delete(h.RoundState, gameId)
	}
}

// Retrieves from the database everything EndRound needs:
// 	- the results of every round played so far, for every player still in the game:
// 	the round specified first, then the rounds before it, down to round 1
// 	- whether it was the last round
func (h *Hub) roundResults(ctx context.Context, gameId, roundNb int64) ([]ws_handlers.RoundResult, bool, error) {
	roundID, err := h.GuessWsHandler.ResolveRoundID(ctx, gameId, roundNb)
	if err != nil {
		return nil, false, err
	}
	activePlayers, err := h.LobbyRepo.AllPlayers(ctx, gameId)
	if err != nil {
		return nil, false, err
	}
	roundGuesses, err := h.GuessWsHandler.GetRoundGuesses(ctx, roundID)
	if err != nil {
		return nil, false, err
	}
	results, err := h.GuessWsHandler.BuildRoundResults(ctx, gameId, roundID, roundNb, activePlayers, roundGuesses)
	if err != nil {
		return nil, false, err
	}
	isLastRound, err := h.GuessWsHandler.IsLastRound(ctx, gameId, roundNb)
	if err != nil {
		return nil, false, err
	}
	// this round first (already built above), then the rounds before it, down to round 1 (none after round 1)
	previousRounds, err := h.GuessWsHandler.BuildGameResults(ctx, gameId, roundNb-1, activePlayers)
	if err != nil {
		return nil, false, err
	}
	gameResults := append([]ws_handlers.RoundResult{*results}, previousRounds...)
	return gameResults, isLastRound, nil
}

// Finalises the round roundNb: broadcasts the results of ALL active players (whether they guessed or not), then,
// unless it was the last round:
//   - multiplayer: schedules the automatic start of the next round after a pause of preRoundPause (10sec)
//   - solo: schedules nothing, the game waits for the player's "start_round"
// EndRound is the only way out of a running round. If the database fails, nothing has changed yet (state still
// "playing", nothing sent): it arms a retry instead of giving up, otherwise the game would stay stuck in this round.
func (h *Hub) EndRound(ctx context.Context, gameId, roundNb int64) {
	state, ok := h.RoundState[gameId]
	if !ok || state.Phase != phasePlaying || state.RoundNb != roundNb {
		// only the round currently running can be ended
		// (when the last guess of a round is processed at the same instant the round's timer fires)
		return
	}
	solo := state.Solo
	if timer, ok := h.GameTimers[gameId]; ok {
		timer.Stop()
		delete(h.GameTimers, gameId)
	}

	gameResults, isLastRound, err := h.roundResults(ctx, gameId, roundNb)
	if err == nil && isLastRound {
		// the game is over
		err = h.GameRepo.SetFinished(ctx, gameId)
	}
	if err != nil {
		if state.EndRetries < endRoundMaxRetries {
			// the round stays "playing" and nothing has been sent: the same round_timeout event tries again in a moment
			state.EndRetries++
			log.Printf("error : could not end round %d of game %d (retry %d/%d in %s): %s\n",
				roundNb, gameId, state.EndRetries, endRoundMaxRetries, endRoundRetryDelay, err)
			h.armTimer(ctx, gameId, endRoundRetryDelay, "round_timeout", roundNb)
			return
		}
		// still failing after about a minute: the error won't go away by itself (missing data, database down for good)
		log.Printf("error : could not end round %d of game %d after %d retries, the game is stopped: %s\n",
			roundNb, gameId, endRoundMaxRetries, err)
		h.abortGame(ctx, gameId, roundNb, solo)
		return
	}
	if isLastRound {
		// xp, rank and statistics of the players still connected at the end: recomputed from all their finished
		// games, once, without retry (if it fails now, their next game or the next start of the server corrects it)
		h.updatePlayers(ctx, h.connectedPlayers(gameId))
		// the players who are still disconnected leave the game now, otherwise their lobby would stay open and they
		// couldn't create or join another one (the connected players leave when they close the game)
		h.LobbyWsHandler.CloseLobby(ctx, gameId, h.connectedPlayers(gameId))
	}

	// every round played so far (this one first)
	// also kept in the state for a player who reloads the page
	broadcastPayload := GameResultsBroadcast{Rounds: gameResults}
	// after the last round the state has no phase anymore: the game is over
	next := &gameRoundState{Solo: solo, RoundNb: roundNb, GameResults: gameResults}
	if !isLastRound {
		next.Phase = phasePause
		if !solo {
			// calculated once and used as reference for all players
			next.NextRoundAt = time.Now().Add(preRoundPause)
			nextRoundAt := next.NextRoundAt
			broadcastPayload.NextRoundAt = &nextRoundAt
		}
	}
	h.RoundState[gameId] = next

	if isLastRound {
		h.broadcast(gameId, WsMessage{Action: "game_results", Payload: broadcastPayload})
	} else {
		h.broadcast(gameId, WsMessage{Action: "round_results", Payload: broadcastPayload})
	}

	if isLastRound && len(h.Clients[gameId]) == 0 {
		// the game ended with nobody connected (everyone lost his connection)
		delete(h.RoundState, gameId)
	}
	if isLastRound || solo {
		// solo: no timer, the next round will be started by the player's "start_round"
		return
	}
	h.armTimer(ctx, gameId, time.Until(next.NextRoundAt), "start_round", roundNb+1)
}

func (h *Hub) closeClients(ctx context.Context) {
	var clientsCpy []*ClientHub
	for _, clients := range h.Clients {
		for client := range clients {
			clientsCpy = append(clientsCpy, client)
		}
	}
	for _, client := range clientsCpy {
		h.unregisterClient(client)
		if !h.LobbyWsHandler.LeaveLobby(ctx, client.UserId) {
			log.Printf("error : the client %d can't leave the lobby in the database\n", 
				client.UserId)
			continue
		}
	}
}

// Updates the xp, then the rank, then the statistics of the players that finished the game
func (h *Hub) updatePlayers(ctx context.Context, connected map[int64]bool) {
	userIds := make([]int64, 0, len(connected))
	for userId := range connected {
		userIds = append(userIds, userId)
	}
	if err := h.UserRepo.UpdateXP(ctx, userIds); err != nil {
		log.Printf("error : could not update the xp of the players %v: %s\n", userIds, err)
	}
	if err := h.UserRepo.UpdateRank(ctx, userIds); err != nil {
		log.Printf("error : could not update the rank of the players %v: %s\n", userIds, err)
	}
	if err := h.StatsRepo.UpdateStats(ctx, userIds); err != nil {
		log.Printf("error : could not update the statistics of the players %v: %s\n", userIds, err)
	}
}

func JoinLobby(ctx context.Context, hub *Hub, wg *sync.WaitGroup) gin.HandlerFunc {
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
            // to avoid a client's upgrade request to be left hanging (= avoid a hung connection)
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
		securesend.SecureSend(ctx, hub.Register, clientHub,
			fmt.Sprintf("Chat : Client %d can't be registered", clientHub.UserId))
		wg.Add(1)
		go clientHub.Connx.WritePump(wg)
		actionClient := func(data []byte) {
			var dataInfos DataInfos
			err := json.Unmarshal(data, &dataInfos)
			if err != nil {
				log.Printf("error : %s\n", err)
				return
			}
			dataInfos.Client = clientHub
			securesend.SecureSend(ctx, hub.Action, &dataInfos,
				fmt.Sprintf("Chat : Client %d can't do this action", clientHub.UserId))
		}
		unregisterClient := func() {
			securesend.SecureSend(ctx, hub.Unregister, clientHub,
				fmt.Sprintf("Chat : Client %d can't be unregistered", clientHub.UserId))
		}
		wg.Add(1)
		clientHub.Connx.ReadPump(actionClient, unregisterClient, wg)
	}
}

// dropIfAbandoned stops a game that nobody is part of anymore (every player left it with "leave"): its timer
// would otherwise keep running rounds for nobody, and a solo game would wait forever with its state in memory.
// A disconnected player still counts as part of the game (this player can come back)
func (h *Hub) dropIfAbandoned(ctx context.Context, gameId int64) {
	players, err := h.LobbyRepo.AllPlayers(ctx, gameId)
	if err != nil || len(players) > 0 {
		return
	}
	if timer, ok := h.GameTimers[gameId]; ok {
		timer.Stop()
		delete(h.GameTimers, gameId)
	}
	delete(h.RoundState, gameId)
}

// handleDisconnection decides what a closed connection means for the player, depending on where the game is
func (h *Hub) handleDisconnection(ctx context.Context, client *ClientHub) {
	if h.connectedPlayers(client.GameId)[client.UserId] {
		// the player still has another tab open on this game
		return
	}
	if h.inProgress(client.GameId) {
		// connection lost during the game: the player stays in it and can reconnect
		// The player only leaves with a "leave" action, or at the end of the game (if no comeback in the lobby)
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

// safely runs one event of the hub. Every game of the server runs in this single goroutine: without it, a single
// bug (a nil pointer...) in one game stops the whole server. The panic is logged with its stack, the hub goes on.
func (h *Hub) safely(event string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in the hub (%s): %v\n%s", event, r, debug.Stack())
		}
	}()
	fn()
}

func (h *Hub) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case client := <- h.Register:
			h.safely("register", func() {
				// when a player comes back after losing his connection
				h.registerClient(client)
				h.clientLobbyBroadcast(client.GameId, ctx)
				if data, ok := h.LobbyWsHandler.SendInfos(ctx, client.GameId); ok {
					h.sendToClient(client, WsMessage{Action: "lobby_infos", Payload: data})
				}
				// then where the game is (countdown, running round or results), so a reconnecting player is back in his round
				h.resyncClient(ctx, client)
			})
		case client := <-h.Unregister:
			h.safely("unregister", func() {
				h.unregisterClient(client)
				h.handleDisconnection(ctx, client)
				// a finished game keeps its state while players are still on its results,
				// and loses it when the last one is gone: nobody can come back to a finished game
				if len(h.Clients[client.GameId]) == 0 && h.gameOver(client.GameId) {
					delete(h.RoundState, client.GameId)
				}
			})
		case serializedData := <- h.Action:
			h.safely("action "+serializedData.Action, func() {
				h.handleAction(ctx, serializedData)
			})
		case <- ctx.Done():
			cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
			defer cancel()
			h.closeClients(cleanupCtx)
			log.Println("The program is finished")
			return
		}
	}
} 