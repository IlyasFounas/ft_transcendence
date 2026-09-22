package ws_handlers

import (
	"backend/internal/models/guesses"
	"backend/internal/models/rounds"
	"backend/internal/repositories"
	// "backend/internal/websockets/handlers"
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"
)

type RoundWsHandler struct {
	roundRepo		*repositories.RoundRepository
	lobbyRepo		*repositories.LobbyRepository
	locationRepo	*repositories.LocationRepository
	GuessWsHandler	*GuessWsHandler

	RoundTimers		map[int64]*time.Timer // means map[gameId]timer
	RoundState		map[int64]*gameRoundState // means map[gameId]currentState, used to resync reconnecting clients
}

func NewRoundWsHandler(
	roundRepo		*repositories.RoundRepository,
	lobbyRepo		*repositories.LobbyRepository,
	locationRepo	*repositories.LocationRepository,
	GuessWsHandler	*GuessWsHandler,
	
) *RoundWsHandler {
	return &RoundWsHandler{
		roundRepo:		roundRepo,
		lobbyRepo:		lobbyRepo,
		locationRepo:	locationRepo,
		GuessWsHandler:	GuessWsHandler,
		RoundTimers:	make(map[int64]*time.Timer),
		RoundState:		make(map[int64]*gameRoundState),
	}
}

// encapsulation of the message sent to the client, with a type and a data payload interface{}
// that can be any struct (RoundStartedBroadcast, RoundResultBroadcast, etc.)
type WsMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// duration of the pause between the end of a round and the start of the next one (if any), no pause if it was the last round.
const postRoundPause = 10 * time.Second

// roundPhase describes what h.roundState[gameId] currently represents
// so a resyncing client gets the right kind of message (a running round vs the pause between rounds)
type roundPhase string

const (
	phasePlaying	roundPhase = "playing" // a round is currently running: look at RoundEndsAt
	phasePause		roundPhase = "pause"   // between two rounds (or after the last one): look at LastResults/NextRoundAt
)

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

func (h *RoundWsHandler) UnmarshalInternalRoundEvent(payload json.RawMessage) (internalRoundEvent, error) {
	var evt internalRoundEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		log.Printf("error : %s\n", err)
		return internalRoundEvent{}, err
	}
	return evt, nil
}

func (h *RoundWsHandler) HandleInternalRoundStart(ctx context.Context, payload json.RawMessage) {
	evt, err := h.UnmarshalInternalRoundEvent(payload)
	if err != nil {
		return
	}
	h.StartRound(ctx, evt.GameId, evt.RoundNb)
}

func (h *RoundWsHandler) HandleInternalRoundTimeout(ctx context.Context, payload json.RawMessage) {
	evt, err := h.UnmarshalInternalRoundEvent(payload)
	if err != nil {
		return
	}
	isLastRound, err := h.isLastRound(ctx, evt.GameId, evt.RoundNb)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	h.EndRound(ctx, evt.GameId, evt.RoundNb, true, isLastRound)
}

func (h *RoundWsHandler) DeleteRoundState(gameId int64) {
	if state, ok := h.RoundState[gameId]; ok && state.IsLastRound {
		delete(h.RoundState, gameId)
	}
}

// this starts the countdown for the round's guess phase: when the timer expires,
// the round is ended even if not all players have guessed
func (h *RoundWsHandler) armRoundTimer(gameId, roundNb int64, duration time.Duration) {
	if existing, ok := h.RoundTimers[gameId]; ok {
		existing.Stop()
	}
	h.RoundTimers[gameId] = time.AfterFunc(duration, func() {
		payload, err := json.Marshal(internalRoundEvent{GameId: gameId, RoundNb: roundNb})
		if err != nil {
			log.Printf("error : %s\n", err)
			return
		}
		h.Action <- &DataInfos{Action: "round_timeout", Payload: payload} //TODO to be put in hub, receives only payload, and sends th rest from hub
	})
}

// this decides, from the hub's state, which round a "start_round" request from a client should actually restart from (but never from anything the client sent):
//   - no state at all: the game hasn't started, so it's round 1
//   - phasePlaying: a round is already running, so this is just a client resyncing
//     (StartRound/EnsureRound is idempotent for a round that already exists, which means it won't recreate it or extend its timer)
//   - phasePause: we're between two rounds, so it's the next one
func (h *RoundWsHandler) ResolveRoundToStart(gameId int64) int64 {
	state, ok := h.RoundState[gameId]
	if !ok {
		return 1
	}
	if state.Phase == phasePlaying || state.IsLastRound {
		return state.RoundNb
	}
	return state.RoundNb + 1
}

// Retrieves the internal ID of a round from its round_nb: the front never needs to know round_id
func (h *RoundWsHandler) resolveRoundID(ctx context.Context, gameID, roundNb int64) (int64, error) {
	return h.roundRepo.GetIDByGameIDAndRoundNb(ctx, gameID, roundNb)
}

// Checks if roundNb is the last round of the game gameID.
func (h *RoundWsHandler) isLastRound(ctx context.Context, gameID, roundNb int64) (bool, error) {
	return h.roundRepo.CheckLastRound(ctx, gameID, roundNb)
}

func (h *RoundWsHandler) GetRoundGuesses(ctx context.Context, roundID int64) []*guesses.Guess {
	roundGuesses, err := h.GuessWsHandler.GetRoundGuesses(ctx, roundID)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil
	}
	return roundGuesses
}

// Finalises the round roundNb: broadcasts the results of ALL active players (whether they guessed or not),
// then schedules a 10s pause before the next round — unless it was the last round.
func (h *RoundWsHandler) EndRound(ctx context.Context, gameId, roundNb int64, timeout , isLastRound bool) (WsMessage, time.Time, error) {
	if state, ok := h.RoundState[gameId]; ok && state.Phase == phasePause && state.RoundNb == roundNb {
		// happens when the last guess of a round is processed at almost the same instant the round's timer fires
		return WsMessage{}, time.Time{}, errors.New("round already ended")
	}
	roundID, err := h.resolveRoundID(ctx, gameId, roundNb)
	if err != nil {
		return WsMessage{}, time.Time{}, err
	}
	roundGuesses := h.GetRoundGuesses(ctx, roundID)
	if (roundGuesses == nil) {
		return WsMessage{}, time.Time{}, errors.New("could not retrieve round guesses")
	}
	activePlayers, err := h.lobbyRepo.AllPlayers(ctx, gameId)
	if err != nil {
		return WsMessage{}, time.Time{}, err
	}
	// checks if we can end the round (all players must have guessed)
	if len(roundGuesses) < len(activePlayers) && !timeout {
		return WsMessage{}, time.Time{}, errors.New("not all players have guessed yet") // not all players have guessed yet, so we cannot end the round
	}
	// stop the round timer if it was still running (it may have already fired, in which case it's a no-op)
	if (!timeout) {
		if timer, ok := h.RoundTimers[gameId]; ok {
			timer.Stop()
			delete(h.RoundTimers, gameId)
		}
	}
	results, err := h.GuessWsHandler.BuildRoundResults(ctx, gameId, roundID, roundNb, activePlayers, roundGuesses)
	if err != nil {
		return WsMessage{}, time.Time{}, err
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

	return WsMessage{Type: "round_results", Data: broadcastPayload}, nextRoundAt, nil
	// h.broadcast(gameId, WsMessage{Type: "round_results", Data: broadcastPayload})

	// if !isLastRound {
	// 	nextRoundNb := roundNb + 1
	// 	time.AfterFunc(time.Until(nextRoundAt), func() {
	// 		payload, err := json.Marshal(internalRoundEvent{GameId: gameId, RoundNb: nextRoundNb})
	// 		if err != nil {
	// 			log.Printf("error : %s\n", err)
	// 			return
	// 		}
	// 		h.Action <- &DataInfos{Action: "start_round", Payload: payload}
	// 	})
	// }
}

func (h *RoundWsHandler) StartAutoNextRound(ctx context.Context, gameId, roundNb, nextRoundNb int64) json {
	payload, err := json.Marshal(internalRoundEvent{GameId: gameId, RoundNb: nextRoundNb})
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil
	}
	return payload
}

// resyncClient sends the newly reconnected client whatever the hub currently knows about the game's round state
func (h *RoundWsHandler) GetRoundState(gameId int64) (bool, WsMessage) {
	state, ok := h.RoundState[gameId]
	if !ok {
		// no round has started yet for this game (e.g. still in the lobby) : nothing to resync.
		return false, WsMessage{}
	}
	switch state.Phase {
	case phasePlaying:
		return true, WsMessage{
			Type: "round_started",
			Data: RoundStartedBroadcast{
				RoundNb: state.RoundNb, 
				RoundEndsAt: state.RoundEndsAt.UTC(),
				MID: state.MID},
		}
	case phasePause:
		payload := RoundResultBroadcast{Results: state.LastResults, IsLastRound: state.IsLastRound}
		if !state.IsLastRound {
			nextRoundAt := state.NextRoundAt
			payload.NextRoundAt = &nextRoundAt
		} else {
			//set finished_at in the database for the game, since the last round has ended
			err := h.GameRepo.SetFinished(context.Background(), gameId)
			if err != nil {
				log.Printf("error : %s\n", err)
				return false, WsMessage{}
			}
		}
		return true, WsMessage{Type: "round_results", Data: payload}
	}
}

// Finds the round for gameID and roundNb, or creates it if it doesn't exist yet. 
// Returns the full round (not just its ID) because we need to retrieve StartedAt as well.
func (h *RoundWsHandler) EnsureRound(ctx context.Context, gameID, roundNb int64) (*rounds.Round, error) {
	// check if the round already exists
	roundID, err := h.roundRepo.GetIDByGameIDAndRoundNb(ctx, gameID, roundNb)
	// if it exists, return it because we need to retrieve StartedAt as well
	if err == nil {
		return h.roundRepo.GetByID(ctx, roundID)
	}
	// if any other error than "round not found" occurs, return it
	if !errors.Is(err, repositories.ErrRoundNotFound) {
		return nil, err
	}

	// TODO only host can create the round
	// TODO if multiplayer, only first round can be created by the host, the rest are created automatically
	// else, host can create each round in order, but cannot create a round if the previous one is not finished yet (or if the game is finished)

	//get the locationID for the round, to create the round
	locationID, err := h.roundRepo.GetLocationID(ctx, gameID, roundNb)
	if err != nil {
		return nil, err
	}

	// check if this is the last round of the game
	isLastRound, err := h.roundRepo.CheckLastRound(ctx, gameID, roundNb)
	if err != nil {
		return nil, err
	}

	return h.roundRepo.Create(ctx, &rounds.Round{
		GameId:     gameID,
		RoundNb:    roundNb,
		LocationId: locationID,
		LastRound:  isLastRound,
		StartedAt:  time.Now(),
	})
}

// It arms the round timer for the given gameId and roundNb, and broadcasts the round_started event to all clients.
// The timer is set to the remaining duration of the round, calculated from the round's StartedAt timestamp and the round duration.
// If the round has already expired (remaining <= 0), it immediately ends the round instead of arming a negative timer.
// This ensures that clients who reconnect late or experience double triggers do not get more time than expected, and keeps all clients
// synchronized with the actual round start time.
func (h *RoundWsHandler) StartRound(ctx context.Context, gameId int64, roundNb int64) {
	// TODO check if the player is authorized to start a round
	// Ensure the round exists and get its duration
	round, err := h.EnsureRound(ctx, gameId, roundNb)
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
		h.EndRound(ctx, gameId, roundNb, false, round.LastRound)
		return
	}

	MID, err := h.locationRepo.GetMIDByID(ctx, round.LocationId)
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