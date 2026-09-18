package ws_handlers

import (
	"backend/internal/handlers"
	"backend/internal/models/games"
	"backend/internal/models/guesses"
	"backend/internal/models/rounds"
	"backend/internal/repositories"
	"backend/internal/tools/haversine"
	"context"
	"encoding/json"
	"errors"
	"github.com/sams96/rgeo"
	"log"
	"time"
)

type GuessWsHandler struct {
	guessRepo *repositories.GuessRepository
	roundRepo *repositories.RoundRepository
	lobbyRepo *repositories.LobbyRepository
	gameRepo  *repositories.GameRepository
	userRepo  *repositories.UserRepository
	reverser  *rgeo.Rgeo
}

func NewGuessWsHandler(
	guessRepo *repositories.GuessRepository,
	roundRepo *repositories.RoundRepository,
	lobbyRepo *repositories.LobbyRepository,
	gameRepo *repositories.GameRepository,
	userRepo *repositories.UserRepository,
	reverser *rgeo.Rgeo,
) *GuessWsHandler {
	return &GuessWsHandler{
		guessRepo: guessRepo,
		roundRepo: roundRepo,
		lobbyRepo: lobbyRepo,
		gameRepo:  gameRepo,
		userRepo:  userRepo,
		reverser:  reverser,
	}
}

// Groups the info about a location (lat, long, country code) used for both the actual location and the guessed point
// useful in case of reco and the front loses the actual location guessed (then both are resent to the front for display)
type LocationInfo struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
	CC   string  `json:"country_code"`
}

// Struct containing all the info about a round, sent through PlayerRoundResult
type RoundInfos struct {
	Distance        int64        `json:"distance"`
	Score           int64        `json:"score"`
	BaseScore       int64        `json:"base_score"`
	CountryFound    bool         `json:"country_found"`
	Base            LocationInfo `json:"base"`
	Guess           LocationInfo `json:"guess"`
	Streak          int64        `json:"streak"`
	TimeLeftSeconds int64        `json:"time_left_seconds"`
}

// Struct containing the result of a round for a player, sent to the front for display and scoring
// sent to the front for each player in the lobby, even if they did not guess 
// (for their own HasGuessed:false and RoundInfos nil)
// TODO missing MID ??
type PlayerRoundResult struct {
	UserID     int64       `json:"user_id"`
	Username   string      `json:"username"`
	RoundNb    int64       `json:"round_nb"`
	HasGuessed bool        `json:"has_guessed"`
	RoundInfos *RoundInfos `json:"round_infos,omitempty"`
	TotalScore int64       `json:"total_score"` // partial when we have not reached the last round, final when we have reached the last round
}

// Validates and calculates (distance/cc/score) and persists a player's guess for the current round.
func (h *GuessWsHandler) AddGuess(ctx context.Context, payload json.RawMessage, userID, gameID int64) (*guesses.Guess, bool) {
	var guess guesses.Guess
	if err := json.Unmarshal(payload, &guess); err != nil {
		log.Printf("error : %s\n", err)
		return nil, false
	}
	guess.UserID = userID

	if guess.Lat == nil || guess.Long == nil {
		log.Println("error : latitude and longitude must be provided")
		return nil, false
	}
	if guess.RoundNb <= 0 {
		log.Println("error : round_nb must be provided")
		return nil, false
	}
	if guess.UserID <= 0 {
		log.Println("error : user_id must be provided")
		return nil, false
	}

	// TODO check if the player has not left the lobby (h.lobbyRepo.ActualLobby)

	var err error
	guess.RoundID, err = h.roundRepo.GetIDByGameIDAndRoundNb(ctx, gameID, guess.RoundNb)
	if err != nil {
		log.Printf("error getting round id for game %d and round_nb %d: %s\n", gameID, guess.RoundNb, err)
		return nil, false
	}

	finished, err := h.gameRepo.Finished(ctx, gameID)
	if err != nil {
		log.Printf("error checking if game %d is finished: %s\n", gameID, err)
		return nil, false
	}
	if finished != nil && !finished.After(time.Now()){
		log.Printf("game %d is already finished, user %d cannot guess\n", gameID, guess.UserID)
		return nil, false
	}

	nextRoundNb, err := h.roundRepo.GetNextRoundNb(ctx, gameID)
	if err != nil {
		log.Printf("error getting current round number for game %d: %s\n", gameID, err)
		return nil, false
	}
	if guess.RoundNb != nextRoundNb-1 {
		log.Printf("user %d is trying to guess for round %d, but the current round for game %d is %d\n",
			guess.UserID, guess.RoundNb, gameID, nextRoundNb-1)
		return nil, false
	}

	locationID, err := h.guessRepo.GetLocationIDByRoundID(ctx, guess.RoundID)
	if err != nil {
		log.Printf("error getting location id for user %d at round %d: %s\n", guess.UserID, guess.RoundID, err)
		return nil, false
	}

	baseLat, baseLong, err := h.roundRepo.GetCoordinatesByLocationID(ctx, locationID)
	if err != nil {
		log.Printf("error getting base coordinates for round %d: %s\n", guess.RoundID, err)
		return nil, false
	}

	baseCC, err := h.roundRepo.GetCCByLocationID(ctx, locationID)
	if err != nil {
		log.Printf("error getting country code for location %d: %s\n", locationID, err)
		return nil, false
	}

	loc, err := h.reverser.ReverseGeocode([]float64{*guess.Long, *guess.Lat})
	if err != nil {
		log.Printf("error reverse geocoding guess coordinates (%f, %f): %s\n", *guess.Lat, *guess.Long, err)
		return nil, false
	}
	guess.CC = loc.CountryCode2
	guess.CountryFound = guess.CC == baseCC

	guess.Distance = handlers.CalculateDistance(
		haversine.Coordinates{Lat: baseLat, Lon: baseLong},
		haversine.Coordinates{Lat: *guess.Lat, Lon: *guess.Long},
	)

	// Time left is calculated here to avoid discrepancies between players in multiplayer
	// TODO put it all in one function
	round, err := h.roundRepo.GetByID(ctx, guess.RoundID)
	if err != nil {
		log.Printf("error getting round %d: %s\n", guess.RoundID, err)
		return nil, false
	}
	duration, err := h.GetRoundDuration(ctx, gameID)
	if err != nil {
		log.Printf("error getting round duration for game %d: %s\n", gameID, err)
		return nil, false
	}
	guess.TimeLeft = duration - time.Since(round.StartedAt)
	if guess.TimeLeft < 0 {
		guess.TimeLeft = 0
	}

	// retrieval of the last streak (avoids client injection)
	guess.Streak, err = h.guessRepo.GetLastStreakByGameAndUserID(ctx, gameID, guess.UserID, guess.RoundNb)
	if err != nil {
		log.Printf("error getting last streak for user %d in game %d: %s\n", guess.UserID, gameID, err)
		return nil, false
	}

	handlers.CalculateScore(&guess, duration)

	createdGuess, err := h.guessRepo.Create(ctx, &guess)
	if err != nil {
		if errors.Is(err, repositories.ErrGuessNotFound) || err.Error() == "pq: duplicate key value violates unique constraint" {
			log.Printf("user %d has already submitted a guess for round %d\n", guess.UserID, guess.RoundID)
			return nil, false
		}
		log.Printf("error creating guess: %s\n", err)
		return nil, false
	}

	return createdGuess, true
}

// Returns the duration of a round for gameID, defined server-side by the game settings
func (h *GuessWsHandler) GetRoundDuration(ctx context.Context, gameID int64) (time.Duration, error) {
	timerSeconds, err := h.gameRepo.Timer(ctx, gameID)
	if err != nil {
		return 0, err
	}
	return time.Duration(timerSeconds) * time.Second, nil
}

// Retrieves all the guesses for a given roundID
func (h *GuessWsHandler) GetRoundGuesses(ctx context.Context, roundID int64) ([]*guesses.Guess, error) {
	return h.guessRepo.GetByRoundID(ctx, roundID)
}

// Transforms the raw guesses of a round into detailed results per player.
// Builds an entry for EACH active player in the game (activePlayers), 
// not just those who guessed (a player who did not guess in time appears 
// with HasGuessed:false and RoundInfos nil, rather than being absent from the payload).
func (h *GuessWsHandler) BuildRoundResults(
	ctx context.Context,
	gameID, roundID, roundNb int64,
	activePlayers []*games.Lobby,
	roundGuesses []*guesses.Guess,
) ([]PlayerRoundResult, error) {
	locationID, err := h.guessRepo.GetLocationIDByRoundID(ctx, roundID)
	if err != nil {
		return nil, err
	}
	baseLat, baseLong, err := h.roundRepo.GetCoordinatesByLocationID(ctx, locationID)
	if err != nil {
		return nil, err
	}
	baseCC, err := h.roundRepo.GetCCByLocationID(ctx, locationID)
	if err != nil {
		return nil, err
	}

	guessByUser := make(map[int64]*guesses.Guess, len(roundGuesses))
	for _, g := range roundGuesses {
		guessByUser[g.UserID] = g
	}

	results := make([]PlayerRoundResult, 0, len(activePlayers))
	for _, player := range activePlayers {
		user, err := h.userRepo.User(ctx, player.UserId)
		if err != nil {
			log.Printf("error : %s\n", err)
			continue
		}
		totalScore, err := h.guessRepo.GetTotalScoreByGameAndUserID(ctx, gameID, player.UserId)
		if err != nil {
			log.Printf("error : %s\n", err)
			continue
		}

		g, hasGuessed := guessByUser[player.UserId]
		result := PlayerRoundResult{
			UserID:     player.UserId,
			Username:   user.Username,
			RoundNb:    roundNb,
			HasGuessed: hasGuessed,
			TotalScore: totalScore,
		}
		if hasGuessed {
			var guessLat, guessLong float64
			if g.Lat != nil {
				guessLat = *g.Lat
			}
			if g.Long != nil {
				guessLong = *g.Long
			}
			result.RoundInfos = &RoundInfos{
				Distance:        g.Distance,
				Score:           g.Score,
				BaseScore:       g.BaseScore,
				CountryFound:    g.CountryFound,
				Base:            LocationInfo{Lat: baseLat, Long: baseLong, CC: baseCC},
				Guess:           LocationInfo{Lat: guessLat, Long: guessLong, CC: g.CC},
				Streak:          g.Streak,
				TimeLeftSeconds: int64(g.TimeLeft.Seconds()),
			}
		}
		results = append(results, result)
	}
	return results, nil
}

// Retrieves the internal ID of a round from its round_nb: the front never needs to know round_id
func (h *GuessWsHandler) ResolveRoundID(ctx context.Context, gameID, roundNb int64) (int64, error) {
	return h.roundRepo.GetIDByGameIDAndRoundNb(ctx, gameID, roundNb)
}

// Checks if roundNb is the last round of the game gameID.
func (h *GuessWsHandler) IsLastRound(ctx context.Context, gameID, roundNb int64) (bool, error) {
	return h.roundRepo.CheckLastRound(ctx, gameID, roundNb)
}

// Finds the round for gameID and roundNb, or creates it if it doesn't exist yet. 
// Returns the full round (not just its ID) because we need to retrieve StartedAt as well.
func (h *GuessWsHandler) EnsureRound(ctx context.Context, gameID, roundNb int64) (*rounds.Round, error) {
	roundID, err := h.roundRepo.GetIDByGameIDAndRoundNb(ctx, gameID, roundNb)
	if err == nil {
		return h.roundRepo.GetByID(ctx, roundID)
	}
	if !errors.Is(err, repositories.ErrRoundNotFound) {
		return nil, err
	}

	locationID, err := h.roundRepo.GetLocationID(ctx, gameID, roundNb)
	if err != nil {
		return nil, err
	}
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