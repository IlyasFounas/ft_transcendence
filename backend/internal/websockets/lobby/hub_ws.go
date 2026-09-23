// Package lobby_ws is the websocket transport layer: it owns client connections, the per-game
// broadcast fan-out, and the single-goroutine event loop (Hub.Run) that serializes every action
// (client-sent or internally generated) so no two actions for the same game are ever handled
// concurrently. All actual game logic — rounds, guesses, lobbies, chat — lives in package
// ws_handlers; this package only calls into it and moves bytes to/from clients.
package lobby_ws

import (
	"backend/internal/repositories"
	"backend/internal/websockets/client"
	"backend/internal/websockets/handlers"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type kickPlayer struct {
	TargetId int64 `json:"target_id"`
}

type playerInfos struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar"`
	Host     bool    `json:"host"`
}

type ClientHub struct {
	Connx   *client_ws.Connexion
	UserId  int64
	LobbyId int64
	GameId  int64
}

// DataInfos carries one action through the hub, whether it comes from a real client or is
// generated internally by a round timer (see ws_handlers.EnqueueFunc).
type DataInfos struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
	// Client is never set by a client — it's always overwritten with the real connection right
	// after unmarshaling (see JoinLobby's actionClient below), and json:"-" means json.Unmarshal
	// ignores it entirely even if a client's payload happens to include a "Client" key. It's nil
	// for events generated internally by the hub.
	Client *ClientHub `json:"-"`
}

type Hub struct {
	Clients        map[int64]map[*ClientHub]bool
	Register       chan *ClientHub
	Unregister     chan *ClientHub
	Action         chan *DataInfos
	LobbyRepo      *repositories.LobbyRepository
	GameRepo       *repositories.GameRepository
	UserRepo       *repositories.UserRepository
	LobbyWsHandler *ws_handlers.LobbyWsHandler
	RoundWsHandler *ws_handlers.RoundWsHandler
}

func NewHub(
	repo *repositories.Repositories,
	wsHandlers *ws_handlers.WsHandlers,
) *Hub {
	h := &Hub{
		Clients:        make(map[int64]map[*ClientHub]bool),
		Register:       make(chan *ClientHub),
		Unregister:     make(chan *ClientHub),
		Action:         make(chan *DataInfos),
		LobbyRepo:      repo.Lobby,
		GameRepo:       repo.Game,
		UserRepo:       repo.User,
		LobbyWsHandler: wsHandlers.Lobby,
		RoundWsHandler: wsHandlers.Round,
	}
	// RoundWsHandler's round timers live in package ws_handlers, which cannot import this
	// package back (that would be an import cycle), so it can't push onto h.Action directly.
	// This callback is the bridge: whenever a round timer fires, RoundWsHandler calls back in
	// here, and we enqueue the resulting action exactly like a real client message.
	wsHandlers.Round.SetEnqueue(func(action string, payload json.RawMessage) {
		h.Action <- &DataInfos{Action: action, Payload: payload}
	})
	return h
}

func NewClientHub(connx *client_ws.Connexion,
	userId int64, lobbyId int64, gameId int64) *ClientHub {
	return &ClientHub{Connx: connx,
		UserId:  userId,
		LobbyId: lobbyId,
		GameId:  gameId}
}

// marshalOrLog marshals payload to JSON once, logging and reporting failure if it can't — the
// shared first step of every send below, which used to marshal (and handle a marshal error)
// separately in broadcast, sendToClient and sendToClientByUserID.
// -> serialization means that the payload is converted to JSON format before being sent to the clients
// -> deserialization means that the backend will convert the JSON back to the original data structure when it receives it
func marshalOrLog(payload interface{}) ([]byte, bool) {
	serializedData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, false
	}
	return serializedData, true
}

// sendMatching marshals payload once and delivers it to every client of gameId for which match
// returns true. broadcast and sendToClientByUserID are both one-off cases of this: everyone in
// the game, or just whichever client in it currently belongs to a given userId.
func (h *Hub) sendMatching(gameId int64, payload interface{}, match func(*ClientHub) bool) {
	serializedData, ok := marshalOrLog(payload)
	if !ok {
		return
	}
	for client := range h.Clients[gameId] {
		if match(client) {
			client.Connx.Send <- serializedData
		}
	}
}

// broadcast serializes the payload once and sends it to all currently connected clients for gameId.
func (h *Hub) broadcast(gameId int64, payload interface{}) {
	h.sendMatching(gameId, payload, func(*ClientHub) bool { return true })
}

func (h *Hub) sendToClientByUserID(gameId int64, userId int64, payload interface{}) {
	h.sendMatching(gameId, payload, func(c *ClientHub) bool { return c.UserId == userId })
}

// sendToClient serializes the payload once and sends it directly to ONE single client — unlike
// sendMatching, the caller already holds the exact client to reach (used for resync messages that
// only concern the reconnecting client), so there's no need to search h.Clients for it.
func (h *Hub) sendToClient(client *ClientHub, payload interface{}) {
	serializedData, ok := marshalOrLog(payload)
	if !ok {
		return
	}
	client.Connx.Send <- serializedData
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
	// Uses ws_handlers.WsMessage rather than a second, locally-declared copy of the same shape:
	// this file used to define its own WsMessage type, identical field-for-field to the one in
	// ws_handlers — harmless at runtime (JSON marshaling doesn't care about Go type identity),
	// but two types for one wire concept invites them drifting apart. There's exactly one now.
	h.broadcast(gameId, ws_handlers.WsMessage{Type: "lobby_players", Data: usersLobby})
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
		h.RoundWsHandler.DeleteRoundState(client.GameId)
	}
	close(client.Connx.Send)
}

// Routes an action received (from a client or generated internally) to the appropriate handler.
// Handled sequentially (only one goroutine reads h.Action in Run) : no race condition possible
// between two actions for the same gameId, even if they come from different clients.
func (h *Hub) handleAction(ctx context.Context, data *DataInfos) {
	switch data.Action {
	case "leave":
		data.Client.Connx.Conn.Close()

	case "kick_player":
		h.handleKickPlayer(ctx, data)

	case "guess":
		h.handleGuess(ctx, data)

	case "start_game":
		h.handleStartGame(ctx, data)

	case "round_begin":
		h.handleRoundBegin(ctx, data)

	case "round_timeout":
		h.handleRoundTimeout(ctx, data)

	case "finished":
		h.handleGameFinished(ctx, data)

	default:
		log.Printf("unknown action : %s\n", data.Action)
	}
}

// finalizeBroadcast broadcasts msg to gameId's clients, unless err is set. ErrRoundNotReady and
// ErrRoundAlreadyEnded are expected, frequent outcomes of EndRound (not every player has guessed
// yet, or a guess and the round timer raced each other) rather than real failures, so they're
// swallowed silently; any other error is logged.
func (h *Hub) finalizeBroadcast(gameId int64, msg ws_handlers.WsMessage, err error) {
	if err != nil {
		if errors.Is(err, ws_handlers.ErrRoundNotReady) || errors.Is(err, ws_handlers.ErrRoundAlreadyEnded) {
			return
		}
		log.Printf("error handling round for game %d: %s\n", gameId, err)
		return
	}
	h.broadcast(gameId, msg)
}

// handleGuess persists the received guess, then finalizes the round if that guess completed it
// (every active player has now guessed). EndRound itself decides whether that's actually the
// case — see finalizeBroadcast for how a "not ready yet" outcome is told apart from a real error.
func (h *Hub) handleGuess(ctx context.Context, data *DataInfos) {
	createdGuess, ok := h.RoundWsHandler.AddGuess(ctx, data.Payload, data.Client.UserId, data.Client.GameId)
	if !ok {
		// TODO : renvoyer une erreur ciblée uniquement à data.Client
		return
	}
	msg, err := h.RoundWsHandler.EndRound(ctx, data.Client.GameId, createdGuess.RoundNb, false)
	h.finalizeBroadcast(data.Client.GameId, msg, err)
}

func (h *Hub) handleKickPlayer(ctx context.Context, data *DataInfos) {
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
			client.Connx.Conn.Close()
			break
		}
	}
	h.clientLobbyBroadcast(data.Client.GameId, ctx)
}

// handleStartGame is the ONLY client-triggered entry point into the round lifecycle: a player
// (the host) asking for the game to begin. It always concerns round 1 — gameId comes from the
// authenticated connection and RequestGameStart itself never takes a round number from the
// client, closing off the "replay/skip rounds at will" hole a client-supplied round number would
// open. RequestGameStart is idempotent (a second click, or another player also starting, just
// resyncs), so this handler doesn't need to guard against being called more than once itself.
func (h *Hub) handleStartGame(ctx context.Context, data *DataInfos) {
	if data.Client == nil {
		// start_game must only ever come from a real client action
		log.Printf("warning : start_game received with no client, ignored\n")
		return
	}
	// TODO or if the client plays in singleplayer, the game should not require a host to start it
	// TODO only the host should be allowed to start the game
	msg, err := h.RoundWsHandler.RequestGameStart(ctx, data.Client.GameId)
	if err != nil {
		log.Printf("error starting game %d: %s\n", data.Client.GameId, err)
		return
	}
	if err := h.LobbyWsHandler.UpdateStartAt(ctx, data.Client.LobbyId); err != nil {
		return
	}
	h.broadcast(data.Client.GameId, msg)
}

// handleRoundBegin fires once a round's pre-round pause has elapsed (see prePlayPause) —
// generated internally by RoundWsHandler's own timer, exactly like handleRoundTimeout below,
// never by a client directly.
func (h *Hub) handleRoundBegin(ctx context.Context, data *DataInfos) {
	if data.Client != nil {
		log.Printf("warning : client %d tried to trigger round_begin directly, ignored\n", data.Client.UserId)
		return
	}
	h.handleInternalRoundEvent(ctx, data.Payload, h.RoundWsHandler.BeginRound)
}

func (h *Hub) handleRoundTimeout(ctx context.Context, data *DataInfos) {
	if data.Client != nil {
		// round_timeout must only ever be generated internally by RoundWsHandler's own timer
		log.Printf("warning : client %d tried to trigger round_timeout directly, ignored\n", data.Client.UserId)
		return
	}
	h.handleInternalRoundEvent(ctx, data.Payload, func(ctx context.Context, gameId, roundNb int64) (ws_handlers.WsMessage, error) {
		return h.RoundWsHandler.EndRound(ctx, gameId, roundNb, true)
	})
}

// handleInternalRoundEvent unmarshals a hub-generated round event payload (round_timeout, or the
// automatic start of the next round) and runs it against run — BeginRound or EndRound — then
// broadcasts whatever it returns. This is the internal-event counterpart to a client-triggered
// action: these events never carry a *ClientHub, so game_id/round_nb travel inside the payload
// instead, which is why run takes them as plain arguments rather than reading data.Client.
func (h *Hub) handleInternalRoundEvent(ctx context.Context, payload json.RawMessage, run func(ctx context.Context, gameId, roundNb int64) (ws_handlers.WsMessage, error)) {
	evt, err := h.RoundWsHandler.UnmarshalInternalRoundEvent(payload)
	if err != nil {
		return
	}
	msg, err := run(ctx, evt.GameId, evt.RoundNb)
	h.finalizeBroadcast(evt.GameId, msg, err)
}

// Checks if the game is finished, if so, closes the client's connection.
func (h *Hub) handleGameFinished(ctx context.Context, data *DataInfos) {
	if data.Client == nil {
		return
	}
	// GameRepo.FinishedAt (not the old Finished bool) matches games.Game.FinishedAt *time.Time
	// and the same check already used in GuessWsHandler.AddGuess: a game is finished once its
	// FinishedAt timestamp is set and has passed.
	_, err := h.GameRepo.FinishedAt(ctx, data.Client.GameId)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	log.Printf("game %d is not finished\n", data.Client.GameId)
}

func JoinLobby(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		// Checks if the player is in the right lobby
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

// Resyncs a client if at least one round has already started, if the game is still ongoing,
// and if the client is not new in this game
func (h *Hub) resyncClient(ctx context.Context, client *ClientHub) {
	// TODO check if the client is not new > do not resync then
	ok, msg := h.RoundWsHandler.GetRoundState(ctx, client.GameId)
	if ok {
		h.sendToClient(client, msg)
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
			h.clientLobbyBroadcast(client.GameId, ctx)
			// send this client (and only this client) the current round state, if any exists
			// so joining mid-game or reconnecting doesn't leave it blind until the next broadcast
			h.resyncClient(ctx, client)
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
		case serializedData := <-h.Action:
			h.handleAction(ctx, serializedData)
		}
	}
}