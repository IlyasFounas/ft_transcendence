package ws_handlers

import (
	"backend/internal/repositories"
	"github.com/sams96/rgeo"
)

type WsHandlers struct {
	Lobby *LobbyWsHandler
	Round *RoundWsHandler
	Guess *GuessWsHandler
	Chat  *ChatHandler
}

func NewWsHandlers(repo *repositories.Repositories, reverser *rgeo.Rgeo) *WsHandlers {
	// Guess is built first, as a local variable, and RoundWsHandler is built from it: RoundWsHandler
	// holds a *GuessWsHandler and delegates guess-related work to it (AddGuess, BuildRoundResults,
	// GetRoundDuration...), so it needs a fully-constructed instance, not a promise of one.
	//
	// Building both inline as struct-literal field values (Round: NewRoundWsHandler(...), Guess:
	// NewGuessWsHandler(...)) would NOT work: Go evaluates composite-literal field values in the
	// order they're written in the source, so "Round:" would still run before "Guess:" regardless
	// of which one is declared first in the WsHandlers type, and RoundWsHandler would be
	// constructed with a nil GuessWsHandler.
	//
	// GuessWsHandler no longer takes repo.Lobby: it never used it (the "check if the player has
	// not left the lobby" TODO it once carried is RoundWsHandler/LobbyWsHandler's job, not this
	// one's), so the unused field and constructor param are gone rather than wired for nothing.
	guessHandler := NewGuessWsHandler(repo.Guess, repo.Round, repo.Game, repo.User, reverser)

	// NOTE: assumes the aggregate Repositories struct exposes the location repository as
	// repo.Location (matching the naming of repo.Lobby, repo.Game, repo.User, etc.) — this is the
	// one field name here I couldn't verify without repositories.go; adjust it if yours differs.
	roundHandler := NewRoundWsHandler(repo.Round, repo.Lobby, repo.Location, repo.Game, guessHandler)

	return &WsHandlers{
		Lobby: NewLobbyWsHandler(repo.Lobby, repo.Game, repo.User),
		Round: roundHandler,
		Guess: guessHandler,
		Chat:  NewChatHandler(repo.Chat, repo.Friend, repo.User),
	}
}