package lobby_ws

import (
	"backend/internal/repositories"
	"backend/internal/websockets/client"
	"backend/internal/websockets/handlers"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {return true},
}

// duration of the pause between the end of a round and the start of the next one (if any), no pause if it was the last round.
const postRoundPause = 10 * time.Second

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
	RoundNb     int64		`json:"round_nb"`
	RoundEndsAt time.Time	`json:"round_ends_at"`
	MID			int64		`json:"mid"`
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
	phasePlaying	roundPhase = "playing" // a round is currently running: look at RoundEndsAt
	phasePause		roundPhase = "pause"   // between two rounds (or after the last one): look at LastResults/NextRoundAt
)

// to be checked when a client reconnects mid-game - refresh, dropped wifi, joining a game already in progress
type gameRoundState struct {
	Phase		roundPhase
	RoundNb		int64
	RoundEndsAt	time.Time // valid when Phase == Playing
	NextRoundAt	time.Time // valid when Phase == Pause and !IsLastRound
	IsLastRound	bool
	LastResults	[]ws_handlers.PlayerRoundResult // results of the round that just ended, valid when Phase == phasePause
	MID			int64 // the Mapillary ID of the location for the current round, valid when Phase == Playing
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
	RoundTimers		map[int64]*time.Timer // means map[gameId]timer
	RoundState		map[int64]*gameRoundState // means map[gameId]currentState, used to resync reconnecting clients
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
		RoundTimers:		make(map[int64]*time.Timer),
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
		select {
		case client.Connx.Send <- serializedData:
		default:
			log.Printf("client_error : the client %d can't read\n", client.UserId)
		}
	}
}

// sendToClient serializes the payload once and sends it to ONE single client 
// unlike broadcast, it is used for resync messages that only concern the reconnecting client
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
	h.broadcast(gameId, WsMessage{Type: "lobby_players", Data: usersLobby})
}

// resyncClient sends the newly reconnected client whatever the hub currently knows about the game's round state
func (h *Hub) resyncClient(client *ClientHub) {
	state, ok := h.RoundState[client.GameId]
	if !ok {
		// no round has started yet for this game (e.g. still in the lobby) : nothing to resync.
		return
	}
	switch state.Phase {
	case phasePlaying:
		h.sendToClient(client, WsMessage{
			Type: "round_started",
			Data: RoundStartedBroadcast{
				RoundNb: state.RoundNb, 
				RoundEndsAt: state.RoundEndsAt.UTC(),
				MID: state.MID},
		})
	case phasePause:
		payload := RoundResultBroadcast{Results: state.LastResults, IsLastRound: state.IsLastRound}
		if !state.IsLastRound {
			nextRoundAt := state.NextRoundAt
			payload.NextRoundAt = &nextRoundAt
		} else {
			//set finished_at in the database for the game, since the last round has ended
			err := h.GameRepo.SetFinished(context.Background(), client.GameId)
			if err != nil {
				log.Printf("error : %s\n", err)
				return
			}
		}
		h.sendToClient(client, WsMessage{Type: "round_results", Data: payload})
	}
}

// this decides, from the hub's state, which round a "start_round" request from a client should actually restart from (but never from anything the client sent):
//   - no state at all: the game hasn't started, so it's round 1
//   - phasePlaying: a round is already running, so this is just a client resyncing
//     (StartRound/EnsureRound is idempotent for a round that already exists, which means it won't recreate it or extend its timer)
//   - phasePause: we're between two rounds, so it's the next one
func (h *Hub) resolveRoundToStart(gameId int64) int64 {
	state, ok := h.RoundState[gameId]
	if !ok {
		return 1
	}
	if state.Phase == phasePlaying || state.IsLastRound {
		return state.RoundNb
	}
	return state.RoundNb + 1
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
		// states' deletion (pause or playing) 
		if state, ok := h.RoundState[client.GameId]; ok && state.IsLastRound {
			delete(h.RoundState, client.GameId)
		}
	}
	close(client.Connx.Send)
}

// Routes an action received (from a client or generated by the hub) to the appropriate handler.
// Handled sequentially (only one goroutine reads h.action in Run) : no race condition possible between
// two actions for the same gameId, even if they come from different clients.
func (h *Hub) handleAction(ctx context.Context, data *DataInfos) {
	switch data.Action {
	case "leave":
		data.Client.Connx.Conn.Close()

	case "kick_player":
		var targetClient kickPlayer
		if err := json.Unmarshal(data.Payload, &targetClient); err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		if !h.LobbyWsHandler.RemovePlayer(ctx, data.Client.UserId, targetClient.TargetId) {
			return
		}
		for client := range h.Clients[data.Client.GameId] {
			if client.UserId == targetClient.TargetId {
				h.clientLobbyBroadcast(data.Client.GameId, ctx)
				client.Connx.Conn.Close()
				return
			}
		}
		h.clientLobbyBroadcast(data.Client.GameId, ctx)

	case "guess":
		h.handleGuess(ctx, data)

	case "start_round":
		if data.Client != nil {
			// TODO or if the client plays in singleplayer, rounds should not start automatically
			
			// client != nil > real client request
			roundNb := h.resolveRoundToStart(data.Client.GameId)
			h.StartRound(ctx, data.Client.GameId, roundNb)
			err := h.LobbyWsHandler.UpdateStartAt(ctx, data.Client.LobbyId, roundNb)
			if err != nil {
				return
			}
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
		//check if the game is finished, if so, update the left_at timestamp in the database and leave the lobby
		if data.Client != nil {
			isFinished, err := h.GameRepo.Finished(ctx, data.Client.GameId)
			if err != nil {
				log.Printf("error : %s\n", err)
				return
			}
			if isFinished {
				log.Printf("game %d finished, client %d left the lobby\n", data.Client.GameId, data.Client.UserId)
				data.Client.Connx.Conn.Close()
				return
			}
			log.Printf("game %d is not finished\n", data.Client.GameId)
			return
		}
	default:
		log.Printf("unknown action : %s\n", data.Action)
	}
}

// handleGuess persists the received guess, then ends the round if all active players in the game have guessed
func (h *Hub) handleGuess(ctx context.Context, data *DataInfos) {
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

// It arms the round timer for the given gameId and roundNb, and broadcasts the round_started event to all clients.
// The timer is set to the remaining duration of the round, calculated from the round's StartedAt timestamp and the round duration.
// If the round has already expired (remaining <= 0), it immediately ends the round instead of arming a negative timer.
// This ensures that clients who reconnect late or experience double triggers do not get more time than expected, and keeps all clients
// synchronized with the actual round start time.
func (h *Hub) StartRound(ctx context.Context, gameId, roundNb int64) {
	// TODO check if the player is authorized to start a round
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

	roundEndsAt := round.StartedAt.Add(duration)
	remaining := time.Until(roundEndsAt)
	if remaining <= 0 {
		// round has expired (ex: late reconnection): end it immediately
		h.EndRound(ctx, gameId, roundNb)
		return
	}

	MID, err := h.LocationRepo.GetMIDByID(ctx, round.LocationId)
	if err != nil {
		log.Printf("error getting MID: %s\n", err)
		return
	}

	h.RoundState[gameId] = &gameRoundState{
		Phase:       phasePlaying,
		RoundNb:     roundNb,
		RoundEndsAt: roundEndsAt,
		MID:         MID, //TODO should be kept server side if we fetch through the proxy, not sent to clients
	}

	h.armRoundTimer(gameId, roundNb, remaining)

	h.broadcast(gameId, WsMessage{
		Type: "round_started",
		Data: RoundStartedBroadcast{
			RoundNb:		roundNb,
			RoundEndsAt:	roundEndsAt.UTC(),
			MID:			MID, 
		},
	})
}

// this starts the countdown for the round's guess phase: when the timer expires,
// the round is ended even if not all players have guessed
func (h *Hub) armRoundTimer(gameId, roundNb int64, duration time.Duration) {
	if existing, ok := h.RoundTimers[gameId]; ok {
		existing.Stop()
	}
	h.RoundTimers[gameId] = time.AfterFunc(duration, func() {
		payload, err := json.Marshal(internalRoundEvent{GameId: gameId, RoundNb: roundNb})
		if err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		h.Action <- &DataInfos{Action: "round_timeout", Payload: payload}
	})
}

// Finalises the round roundNb: broadcasts the results of ALL active players (whether they guessed or not),
// then schedules a 10s pause before the next round — unless it was the last round.
func (h *Hub) EndRound(ctx context.Context, gameId, roundNb int64) {
	if state, ok := h.RoundState[gameId]; ok && state.Phase == phasePause && state.RoundNb == roundNb {
		// happens when the last guess of a round is processed at almost the same instant the round's timer fires
		return
	}
	if timer, ok := h.RoundTimers[gameId]; ok {
		timer.Stop()
		delete(h.RoundTimers, gameId)
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

	broadcastPayload := RoundResultBroadcast{Results: results, IsLastRound: isLastRound}
	var nextRoundAt time.Time
	if !isLastRound {
		// calculated once: this value is the reference. The AfterFunc below is based on the remaining time
		// until this same instant (time.Until), so the next round starts exactly when all clients have seen their countdown reach zero.
		nextRoundAt = time.Now().Add(postRoundPause)
		broadcastPayload.NextRoundAt = &nextRoundAt
	}
		h.RoundState[gameId] = &gameRoundState{
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
		h.Action <- &DataInfos{Action: "start_round", Payload: payload}
	})
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

func JoinLobby(hub *Hub, wg *sync.WaitGroup) gin.HandlerFunc {
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
			hub.Action <- &dataInfos 
		}
		unregisterClient := func() {
			hub.Unregister <- clientHub
		}
		wg.Add(1)
		clientHub.Connx.ReadPump(actionClient, unregisterClient, wg)
	}
}

func (h *Hub) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case client := <- h.Register:
			h.registerClient(client)
			h.clientLobbyBroadcast(client.GameId, ctx)
			// send this client (and only this client) the current round state, if any exists
			// so joining mid-game or reconnecting doesn't leave it blind until the next broadcast
			h.resyncClient(client) // TODO check if client reconnects only if it was already registered before (h.clients[gameId][client] == true) : if so, skip resyncClient because it already has the current state. Otherwise a reconnecting client will get two resyncs (one from the first registration, one from the second) and see the same round_started broadcast twice.
			data, ok := h.LobbyWsHandler.SendInfos(ctx, client.GameId)
			if !ok {
				continue
			}
			h.sendToClientByUserID(client.GameId, client.UserId, data)
		case client := <-h.Unregister:
			h.unregisterClient(client)
			if !h.LobbyWsHandler.LeaveLobby(ctx, client.UserId) {
				log.Println("error : the client can't leave the lobby in the database")
				continue
			}
			h.clientLobbyBroadcast(client.GameId, ctx)
		case serializedData := <- h.Action:
			h.handleAction(ctx, serializedData)
		case <- ctx.Done():
			cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
			defer cancel()
			h.closeClients(cleanupCtx)
			log.Println("The program is finished")
			return 
		}
	}
} 