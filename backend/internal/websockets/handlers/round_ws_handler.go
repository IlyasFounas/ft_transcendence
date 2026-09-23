// Package ws_handlers holds the game-domain logic for the websocket-driven multiplayer flow:
// lobbies (LobbyWsHandler), guesses/scoring (GuessWsHandler), rounds (RoundWsHandler) and chat
// (ChatHandler). None of these types know about websocket connections or clients — that's the
// transport layer, owned by the Hub in package lobby_ws. lobby_ws imports this package (to call
// into the game logic), never the other way around, so nothing in here may import lobby_ws.
package ws_handlers

import (
	"backend/internal/models/guesses"
	"backend/internal/models/rounds"
	"backend/internal/repositories"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

// RoundWsHandler owns everything about "where a game currently is in its round cycle": scheduling
// and arming the pause before each round, firing round timeouts, resolving which round is
// current, building round results, and resyncing a (re)connecting client. It never touches the
// hub's client list or connections directly — see EnqueueFunc below for why, and why every public
// method here returns a message for the caller to broadcast instead of sending it themselves.
//
// # The round lifecycle, end to end
//
// A game's timeline is a strict sequence for each round N: PAUSE (prePlayPause, fixed, always
// present — including before round 1) -> PLAYING (Game.Timer seconds) -> PAUSE (before round
// N+1) -> ... -> PLAYING (last round) -> finished. There is exactly one code path that ever
// creates a round's DB row and decides when it begins: scheduleRoundStart. It is used identically
// whether we're about to play round 1 (triggered by a client's "start_game") or round N+1
// (triggered automatically by EndRound) — round 1 is not a special case anywhere in this file.
//
// A round's row is written to the DB as soon as its pause begins, with StartedAt already set to
// the future instant the pause will elapse — not once play actually starts. This makes StartedAt
// (compared against Game.Timer, fetched live) sufficient on its own to derive, at any later time,
// whether a round is still pending, currently playing, or already over — see BeginRound and
// recoverRoundState. Nothing about a round's duration is ever stored on the round itself: Game.Timer
// is the single source of truth for how long every round of a game lasts, read fresh every time
// it's needed.
type RoundWsHandler struct {
	roundRepo      *repositories.RoundRepository
	lobbyRepo      *repositories.LobbyRepository
	locationRepo   *repositories.LocationRepository
	gameRepo       *repositories.GameRepository
	guessWsHandler *GuessWsHandler

	pendingTimers map[int64]*time.Timer     // map[gameId]timer — the ONE outstanding timer for that game: either the pre-round pause elapsing, or the current round timing out. Never both: a game is always in exactly one phase at a time.
	roundState    map[int64]*gameRoundState // map[gameId]currentState, used to resync reconnecting clients
	enqueue       EnqueueFunc
}

func NewRoundWsHandler(
	roundRepo *repositories.RoundRepository,
	lobbyRepo *repositories.LobbyRepository,
	locationRepo *repositories.LocationRepository,
	gameRepo *repositories.GameRepository,
	guessWsHandler *GuessWsHandler,
) *RoundWsHandler {
	return &RoundWsHandler{
		roundRepo:      roundRepo,
		lobbyRepo:      lobbyRepo,
		locationRepo:   locationRepo,
		gameRepo:       gameRepo,
		guessWsHandler: guessWsHandler,
		pendingTimers:  make(map[int64]*time.Timer),
		roundState:     make(map[int64]*gameRoundState),
	}
}

// EnqueueFunc lets RoundWsHandler push a server-generated action (round_begin once a pause
// elapses, or round_timeout once a round's guess phase runs out) back onto the hub's own serial
// event loop, the same way a real client message arrives. RoundWsHandler cannot depend on the
// hub's package to do this directly — lobby_ws already imports ws_handlers, so the reverse would
// be an import cycle, and DataInfos/h.Action simply don't exist here. Instead, the hub injects
// this callback once, right after constructing both itself and this handler (see
// lobby_ws.NewHub). Until SetEnqueue is called, a timer firing is a silent no-op.
type EnqueueFunc func(action string, payload json.RawMessage)

// SetEnqueue wires the hub's callback into this handler. NewHub calls this once, immediately
// after constructing both the hub and this handler.
func (h *RoundWsHandler) SetEnqueue(fn EnqueueFunc) {
	h.enqueue = fn
}

// encapsulation of the message sent to the client, with a type and a data payload interface{}
// that can be any struct (RoundStartedBroadcast, RoundResultBroadcast, etc.)
type WsMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// prePlayPause is the fixed pause observed before EVERY round begins — round 1 exactly like any
// other. It gives every client's UI a predictable "get ready" beat (loading the street view,
// showing who's about to play, etc.) and a moment to have fully received round N's results before
// round N+1's countdown starts.
const prePlayPause = 10 * time.Second

// roundPhase describes what h.roundState[gameId] currently represents, so a resyncing client
// gets the right kind of message.
type roundPhase string

const (
	phasePause   roundPhase = "pause"   // waiting for a round to begin (round 1's pre-game pause, or the pause between two rounds): look at NextRoundAt
	phasePlaying roundPhase = "playing" // a round is currently running: look at RoundEndsAt
)

// internalRoundEvent carries the game_id/round_nb for events the HUB ITSELF generates — a pause
// elapsing (round_begin) or a round's timer running out (round_timeout). These events never come
// from a client — they're sent through EnqueueFunc from a timer callback.
type internalRoundEvent struct {
	GameId  int64 `json:"game_id"`
	RoundNb int64 `json:"round_nb"`
}

// Sent once a client asks the game to start (round 1) and reused, unchanged, when resyncing a
// client that reconnects during that same initial pause. The pause between two later rounds
// instead rides inside RoundResultBroadcast.NextRoundAt below, since that message already needs
// to carry the round that just ended — no reason to send two separate messages for one pause.
type RoundPendingBroadcast struct {
	NextRoundNb int64     `json:"next_round_nb"`
	NextRoundAt time.Time `json:"next_round_at"`
}

// It will be sent to all clients when a round starts, with the round number and the absolute end time of the round
// (not a relative duration). Each client will calculate its own countdown using (RoundEndsAt - now), so all clients
// will converge to the same final instant regardless of their network latency at the time of reception.
type RoundStartedBroadcast struct {
	RoundNb     int64     `json:"round_nb"`
	RoundEndsAt time.Time `json:"round_ends_at"`
	MID         int64     `json:"mid"`
}

// It shares the results of the round to all players, even those who did not guess (they will see their own result as "not guessed").
// The results are sorted by player ID for consistency across clients.
// It is shared at the end of a round, with an indication of whether it was the last round or not :
// 	-> if it is not the last round, NextRoundAt is provided as an absolute timestamp for when the next round will start.
// It will be nil or omitted if IsLastRound is true.
type RoundResultBroadcast struct {
	// PlayerRoundResult is defined in guess_ws_handler.go — same package, so no qualifier here.
	Results     []PlayerRoundResult `json:"results"`
	IsLastRound bool                `json:"is_last_round"`
	NextRoundAt *time.Time          `json:"next_round_at,omitempty"`
}

// to be checked when a client reconnects mid-game - refresh, dropped wifi, joining a game already in progress
type gameRoundState struct {
	Phase       roundPhase
	RoundNb     int64
	RoundEndsAt time.Time // valid when Phase == phasePlaying
	NextRoundAt time.Time // valid when Phase == phasePause
	IsLastRound bool      // valid when Phase == phasePause and LastResults != nil (i.e. not the pre-game pause)
	LastResults []PlayerRoundResult // results of the round that just ended, valid when Phase == phasePause and this ISN'T the pre-game pause. nil during the pre-game pause — that's how GetRoundState tells the two pauses apart.
	MID         int64                // the Mapillary ID of the location for the current round, valid when Phase == phasePlaying
}

// UnmarshalInternalRoundEvent decodes a hub-generated event payload. Exported because the hub
// needs it too (to learn which game/round an internal "round_begin"/"round_timeout" action
// refers to, before it can call BeginRound/EndRound) — see lobby_ws's handleInternalRoundEvent.
func (h *RoundWsHandler) UnmarshalInternalRoundEvent(payload json.RawMessage) (internalRoundEvent, error) {
	var evt internalRoundEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		log.Printf("error : %s\n", err)
		return internalRoundEvent{}, err
	}
	return evt, nil
}

// DeleteRoundState drops gameId's round state once every client has left AND the game had
// already reached its last round — called from the hub when its client set for a game empties
// out. State for a game that's still mid-way through is kept, so a reconnecting player can still
// be resynced.
func (h *RoundWsHandler) DeleteRoundState(gameId int64) {
	if state, ok := h.roundState[gameId]; ok && state.IsLastRound {
		delete(h.roundState, gameId)
	}
}

// enqueueEvent marshals evt and pushes it back onto the hub's event loop as `action`. Shared by
// every timer callback below — each one used to open-code this exact marshal-then-enqueue
// sequence inside its own time.AfterFunc closure.
func (h *RoundWsHandler) enqueueEvent(action string, evt internalRoundEvent) {
	payload, err := json.Marshal(evt)
	if err != nil {
		log.Printf("error : %s\n", err)
		return
	}
	if h.enqueue != nil {
		h.enqueue(action, payload)
	}
}

// armTimer is the single place a game's outstanding timer gets (re)armed, whether it's counting
// down a pre-round pause or a round's guess phase. A game only ever has ONE meaningful timer in
// flight at a time (see pendingTimers' doc comment), so arming a new one always stops whatever
// was there before — normally a no-op, since the phase that just ended is responsible for
// stopping its own timer, but harmless and cheap to guarantee here too.
func (h *RoundWsHandler) armTimer(gameId int64, duration time.Duration, fire func()) {
	if existing, ok := h.pendingTimers[gameId]; ok {
		existing.Stop()
	}
	h.pendingTimers[gameId] = time.AfterFunc(duration, fire)
}

// stopTimer cancels and forgets gameId's outstanding timer, if any. A no-op if it already fired
// (exactly the case when this is called because that very timer is what triggered the current
// call) or if there was none.
func (h *RoundWsHandler) stopTimer(gameId int64) {
	if timer, ok := h.pendingTimers[gameId]; ok {
		timer.Stop()
		delete(h.pendingTimers, gameId)
	}
}

// Retrieves the internal ID of a round from its round_nb: the front never needs to know round_id
func (h *RoundWsHandler) resolveRoundID(ctx context.Context, gameID, roundNb int64) (int64, error) {
	return h.roundRepo.GetIDByGameIDAndRoundNb(ctx, gameID, roundNb)
}

// Checks if roundNb is the last round of the game gameID.
func (h *RoundWsHandler) isLastRound(ctx context.Context, gameID, roundNb int64) (bool, error) {
	return h.roundRepo.CheckLastRound(ctx, gameID, roundNb)
}

// AddGuess delegates to GuessWsHandler. Exposed here so the hub only ever talks to
// RoundWsHandler for anything guess/round related, never to GuessWsHandler directly.
func (h *RoundWsHandler) AddGuess(ctx context.Context, payload json.RawMessage, userID, gameID int64) (*guesses.Guess, bool) {
	return h.guessWsHandler.AddGuess(ctx, payload, userID, gameID)
}

// ErrRoundNotReady is returned by EndRound when it's called for a non-timeout reason (typically:
// a player just guessed) but not every active player has guessed yet. This is an expected,
// frequent outcome — it happens on every guess except the round's very last one — not a real
// failure, so callers should treat it as a no-op rather than log it as an error.
var ErrRoundNotReady = errors.New("not all players have guessed yet")

// ErrRoundAlreadyEnded is returned when EndRound is called for a round that has already been
// finalized. This happens when the last guess of a round is processed at almost the same instant
// the round's timer fires: both paths independently decide the round is over and each try to end
// it. Like ErrRoundNotReady, callers should treat this as a no-op, not a real failure.
var ErrRoundAlreadyEnded = errors.New("round already ended")

// ensureRoundRow finds roundNb's DB row for gameID, or creates it — with StartedAt set to
// startsAt, the instant its pause is scheduled to elapse — if it doesn't exist yet. It is the
// ONLY place a round row is ever created; called exclusively from scheduleRoundStart, right when
// that round's pause begins.
func (h *RoundWsHandler) ensureRoundRow(ctx context.Context, gameID, roundNb int64, startsAt time.Time) (*rounds.Round, error) {
	roundID, err := h.resolveRoundID(ctx, gameID, roundNb)
	if err == nil {
		return h.roundRepo.GetByID(ctx, roundID)
	}
	if !errors.Is(err, repositories.ErrRoundNotFound) {
		return nil, err
	}

	// TODO only host can trigger the start of round 1
	// TODO if multiplayer, only the host's client should ever cause a round to be scheduled

	locationID, err := h.roundRepo.GetLocationID(ctx, gameID, roundNb)
	if err != nil {
		return nil, err
	}
	isLastRound, err := h.isLastRound(ctx, gameID, roundNb)
	if err != nil {
		return nil, err
	}
	return h.roundRepo.Create(ctx, &rounds.Round{
		GameId:     gameID,
		RoundNb:    roundNb,
		LocationId: locationID,
		LastRound:  isLastRound,
		StartedAt:  startsAt,
	})
}

// scheduleRoundStart is THE single code path that starts a round's pre-round pause: it persists
// roundNb's row (StartedAt = now + prePlayPause) and arms the timer that will fire "round_begin"
// once that pause elapses. Used identically for round 1 (from RequestGameStart) and for round N+1
// (from EndRound) — there is no special-cased "first round" logic anywhere in this file.
//
// It does NOT touch h.roundState — callers set that themselves, because they need to attach
// different extra context to it (round 1 has no previous results; round N+1's pause message also
// carries the round that just ended).
func (h *RoundWsHandler) scheduleRoundStart(ctx context.Context, gameId, roundNb int64) (time.Time, error) {
	startsAt := time.Now().Add(prePlayPause)
	if _, err := h.ensureRoundRow(ctx, gameId, roundNb, startsAt); err != nil {
		return time.Time{}, err
	}
	h.armTimer(gameId, prePlayPause, func() {
		h.enqueueEvent("round_begin", internalRoundEvent{GameId: gameId, RoundNb: roundNb})
	})
	return startsAt, nil
}

// RequestGameStart is the ONLY client-triggered entry point into the round lifecycle: a player
// (the host) asking for the game to begin. It always concerns round 1 — a client can never name
// an arbitrary round to jump to, closing off the "replay/skip rounds at will" hole the old
// client-supplied round number had.
//
// It is idempotent two different ways:
//   - if this game already has in-memory state, the game has already been started (by this
//     client's own double click, or by another player) — this just resyncs instead of scheduling
//     a second, conflicting round 1.
//   - if there is no in-memory state but round 1 already exists in the DB, the process most
//     likely restarted mid-game: recoverRoundState rebuilds the right in-memory state (and
//     re-arms whatever timer is still needed) from that row instead of blindly recreating round 1
//     and resetting the game.
func (h *RoundWsHandler) RequestGameStart(ctx context.Context, gameId int64) (WsMessage, error) {
	if _, ok := h.roundState[gameId]; ok {
		_, msg := h.GetRoundState(ctx, gameId)
		return msg, nil
	}

	if roundID, err := h.resolveRoundID(ctx, gameId, 1); err == nil {
		return h.recoverRoundState(ctx, gameId, roundID)
	} else if !errors.Is(err, repositories.ErrRoundNotFound) {
		return WsMessage{}, err
	}

	startsAt, err := h.scheduleRoundStart(ctx, gameId, 1)
	if err != nil {
		return WsMessage{}, fmt.Errorf("error starting game: %w", err)
	}
	h.roundState[gameId] = &gameRoundState{
		Phase:       phasePause,
		RoundNb:     1,
		NextRoundAt: startsAt,
	}
	return WsMessage{Type: "round_pending", Data: RoundPendingBroadcast{NextRoundNb: 1, NextRoundAt: startsAt.UTC()}}, nil
}

// recoverRoundState rebuilds gameId's in-memory state from round #1's row alone (roundID is
// already resolved by the caller), deriving which phase the game is actually in — still paused,
// mid-round, or its round timer already elapsed while nothing was listening (e.g. the process was
// down) — purely from StartedAt compared to now, then arms whatever timer that phase still needs.
// This is what makes RequestGameStart safe to call again after a restart.
func (h *RoundWsHandler) recoverRoundState(ctx context.Context, gameId, roundID int64) (WsMessage, error) {
	round, err := h.roundRepo.GetByID(ctx, roundID)
	if err != nil {
		return WsMessage{}, err
	}
	if time.Now().Before(round.StartedAt) {
		h.roundState[gameId] = &gameRoundState{Phase: phasePause, RoundNb: round.RoundNb, NextRoundAt: round.StartedAt}
		h.armTimer(gameId, time.Until(round.StartedAt), func() {
			h.enqueueEvent("round_begin", internalRoundEvent{GameId: gameId, RoundNb: round.RoundNb})
		})
		return WsMessage{Type: "round_pending", Data: RoundPendingBroadcast{NextRoundNb: round.RoundNb, NextRoundAt: round.StartedAt.UTC()}}, nil
	}
	// The pause is already over: BeginRound re-derives everything it needs (duration read live
	// from Game, time remaining, MID...) on its own, and already knows to fall back to an
	// immediate EndRound if even the round's own timer has elapsed by now — exactly the "very
	// late reconnection" case it has to handle regardless of why we're calling it.
	return h.BeginRound(ctx, gameId, round.RoundNb)
}

// BeginRound is called ONLY once a round's pre-round pause has elapsed — either just now (the
// normal "round_begin" internal event, fired by the timer scheduleRoundStart armed) or while
// recovering state after a restart. It is never reachable from a client action directly, so
// roundNb here is always trustworthy. It arms this round's own guess-phase timer and returns the
// round_started broadcast for the caller to send.
func (h *RoundWsHandler) BeginRound(ctx context.Context, gameId, roundNb int64) (WsMessage, error) {
	roundID, err := h.resolveRoundID(ctx, gameId, roundNb)
	if err != nil {
		return WsMessage{}, fmt.Errorf("error beginning round: %w", err)
	}
	round, err := h.roundRepo.GetByID(ctx, roundID)
	if err != nil {
		return WsMessage{}, fmt.Errorf("error beginning round: %w", err)
	}

	// Read live from Game every time — a round's duration is NEVER persisted on the round
	// itself. Game.Timer is the single source of truth for how long every round of this game
	// lasts.
	duration, err := h.guessWsHandler.GetRoundDuration(ctx, gameId)
	if err != nil {
		return WsMessage{}, fmt.Errorf("error getting round duration: %w", err)
	}

	roundEndsAt := round.StartedAt.Add(duration)
	remaining := time.Until(roundEndsAt)
	if remaining <= 0 {
		// The round's own duration has already elapsed by the time we got to actually begin it
		// (e.g. the process was down through the whole pause AND the round itself): force-end it
		// immediately (timeout=true) rather than arming a negative timer or waiting on guesses
		// that can no longer come in.
		return h.EndRound(ctx, gameId, roundNb, true)
	}

	MID, err := h.locationRepo.GetMIDByID(ctx, round.LocationId)
	if err != nil {
		return WsMessage{}, fmt.Errorf("error getting MID: %w", err)
	}

	h.roundState[gameId] = &gameRoundState{
		Phase:       phasePlaying,
		RoundNb:     roundNb,
		RoundEndsAt: roundEndsAt,
		MID:         MID, //TODO should be kept server side if we fetch through the proxy, not sent to clients
	}
	h.armTimer(gameId, remaining, func() {
		h.enqueueEvent("round_timeout", internalRoundEvent{GameId: gameId, RoundNb: roundNb})
	})

	return roundStartedMessage(roundNb, roundEndsAt, MID), nil
}

// roundStartedMessage builds the round_started broadcast.
func roundStartedMessage(roundNb int64, roundEndsAt time.Time, mid int64) WsMessage {
	return WsMessage{
		Type: "round_started",
		Data: RoundStartedBroadcast{
			RoundNb:     roundNb,
			RoundEndsAt: roundEndsAt.UTC(),
			MID:         mid,
		},
	}
}

// Finalises the round roundNb: builds the results of ALL active players (whether they guessed or
// not) and, unless it was the last round, schedules round N+1's pre-round pause via
// scheduleRoundStart. timeout distinguishes a natural end (every active player has guessed) from
// a forced one (the round's timer fired before everyone guessed): only a forced end may proceed
// while players are still missing a guess.
func (h *RoundWsHandler) EndRound(ctx context.Context, gameId, roundNb int64, timeout bool) (WsMessage, error) {
	if state, ok := h.roundState[gameId]; ok && state.Phase == phasePause && state.RoundNb == roundNb {
		// happens when the last guess of a round is processed at almost the same instant the round's timer fires
		return WsMessage{}, ErrRoundAlreadyEnded
	}

	roundID, err := h.resolveRoundID(ctx, gameId, roundNb)
	if err != nil {
		return WsMessage{}, err
	}
	// LastRound is read off the row rather than recomputed via isLastRound()/CheckLastRound:
	// ensureRoundRow already decided and persisted it once, when the round's pause began, so
	// asking the same round_nb/round_total question again here would just be re-checking
	// something the code already checked.
	round, err := h.roundRepo.GetByID(ctx, roundID)
	if err != nil {
		return WsMessage{}, err
	}
	roundGuesses, err := h.guessWsHandler.GetRoundGuesses(ctx, roundID)
	if err != nil {
		return WsMessage{}, err
	}
	activePlayers, err := h.lobbyRepo.AllPlayers(ctx, gameId)
	if err != nil {
		return WsMessage{}, err
	}
	// checks if we can end the round (all players must have guessed, unless the timer forced it)
	if len(roundGuesses) < len(activePlayers) && !timeout {
		return WsMessage{}, ErrRoundNotReady
	}

	// stop the round's guess-phase timer — a no-op if it already fired, which is exactly the
	// timeout==true case
	h.stopTimer(gameId)

	results, err := h.guessWsHandler.BuildRoundResults(ctx, gameId, roundID, roundNb, activePlayers, roundGuesses)
	if err != nil {
		return WsMessage{}, err
	}

	broadcastPayload := RoundResultBroadcast{Results: results, IsLastRound: round.LastRound}
	state := &gameRoundState{
		Phase:       phasePause,
		RoundNb:     roundNb,
		IsLastRound: round.LastRound,
		LastResults: results,
	}

	if !round.LastRound {
		nextRoundAt, err := h.scheduleRoundStart(ctx, gameId, roundNb+1)
		if err != nil {
			return WsMessage{}, fmt.Errorf("error scheduling round %d: %w", roundNb+1, err)
		}
		nextRoundAtUTC := nextRoundAt.UTC()
		broadcastPayload.NextRoundAt = &nextRoundAtUTC
		state.NextRoundAt = nextRoundAt
	}

	h.roundState[gameId] = state

	return WsMessage{Type: "round_results", Data: broadcastPayload}, nil
}

// GetRoundState returns whatever the hub currently knows about gameId's round state, for
// resyncing a client that just (re)connected. The bool return is false when there is nothing to
// resync (no round has started yet for this game — i.e. RequestGameStart has never been called).
func (h *RoundWsHandler) GetRoundState(ctx context.Context, gameId int64) (bool, WsMessage) {
	state, ok := h.roundState[gameId]
	if !ok {
		return false, WsMessage{}
	}
	switch state.Phase {
	case phasePlaying:
		return true, roundStartedMessage(state.RoundNb, state.RoundEndsAt, state.MID)
	case phasePause:
		if state.LastResults == nil {
			// the pre-game pause (before round 1): nothing has been played yet, so there's
			// nothing to resync but the countdown to round 1.
			return true, WsMessage{Type: "round_pending", Data: RoundPendingBroadcast{NextRoundNb: state.RoundNb, NextRoundAt: state.NextRoundAt.UTC()}}
		}
		payload := RoundResultBroadcast{Results: state.LastResults, IsLastRound: state.IsLastRound}
		if !state.IsLastRound {
			nextRoundAt := state.NextRoundAt.UTC()
			payload.NextRoundAt = &nextRoundAt
		} else {
			// the last round just ended: record finished_at now that we know for certain no
			// further round will start for this game. The resync still goes out even if this
			// write fails — the client shouldn't get nothing just because this bookkeeping
			// call had a hiccup.
			if err := h.gameRepo.SetFinished(ctx, gameId); err != nil {
				log.Printf("error : %s\n", err)
			}
		}
		return true, WsMessage{Type: "round_results", Data: payload}
	default:
		// unreachable in practice (roundPhase only ever holds phasePlaying/phasePause), but the
		// compiler can't know that from a plain string type, so every path needs a return.
		return false, WsMessage{}
	}
}