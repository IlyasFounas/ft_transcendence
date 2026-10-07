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
	"log"
	"math"
	"time"

	"backend/internal/tools/rgeo"
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
	Guess           LocationInfo `json:"guess"`
	Streak          int64        `json:"streak"`
	TimeLeftSeconds int64        `json:"time_left_seconds"`
}

// Struct containing the result of a round for a player, sent to the front for display and scoring
// sent to the front for each player in the lobby, even if they did not guess
// (for their own HasGuessed:false and RoundInfos nil)
type PlayerRoundResult struct {
	UserID     int64       `json:"user_id"`
	Username   string      `json:"username"`
	HasGuessed bool        `json:"has_guessed"`
	RoundInfos *RoundInfos `json:"round_infos,omitempty"`
	TotalScore int64       `json:"total_score"` // partial when we have not reached the last round, final when we have reached the last round
}

// Results of one round, sent to the front: the round number and the actual location (base) come once, first,
// and are always present (even if nobody guessed), then the result of each player (always present too, possibly empty)
type RoundResult struct {
	RoundNb           int64               `json:"round_nb"`
	Base              LocationInfo        `json:"base"`
	PlayerRoundResult []PlayerRoundResult `json:"results"`
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
	if *guess.Lat > 90 || *guess.Lat < -90 {
		*guess.Lat = math.Mod(*guess.Lat+90, 180) // works between 0 and 180
		if *guess.Lat < 0 { //set back to normal coordinates
			*guess.Lat += 90
		} else {
			*guess.Lat -= 90
		}
	}
	if *guess.Long > 180 || *guess.Long < -180 {
		*guess.Long = math.Mod(*guess.Long+180, 360) // works between 0 and 360
		if *guess.Long < 0 { //set back to normal coordinates
			*guess.Long += 180
		} else {
			*guess.Long -= 180
		}
	}
	if guess.RoundNb <= 0 {
		log.Println("error : round_nb must be provided")
		return nil, false
	}
	if guess.UserID <= 0 {
		log.Println("error : user_id must be provided")
		return nil, false
	}

	var err error
	guess.RoundID, err = h.roundRepo.GetIDByGameIDAndRoundNb(ctx, gameID, guess.RoundNb)
	if err != nil {
		log.Printf("error getting round id for game %d and round_nb %d: %s\n", gameID, guess.RoundNb, err)
		return nil, false
	}

	finished, err := h.gameRepo.FinishedAt(ctx, gameID)
	if err != nil {
		log.Printf("error checking if game %d is finished: %s\n", gameID, err)
		return nil, false
	}
	if finished != nil && !finished.After(time.Now()) {
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
	if err != nil && !errors.Is(err, rgeo.ErrLocationNotFound) {
		log.Printf("error reverse geocoding guess coordinates (%f, %f): %s\n", *guess.Lat, *guess.Long, err)
		return nil, false
	}
	guess.CC = loc.CountryCode2
	guess.CountryFound = guess.CC != "" && guess.CC == baseCC

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

// Transforms the raw guesses of a round into its results
func (h *GuessWsHandler) BuildRoundResults(
	ctx context.Context,
	gameID, roundID, roundNb int64,
	activePlayers []*games.Lobby,
	roundGuesses []*guesses.Guess,
) (*RoundResult, error) {
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

	// round number and actual location once, for everyone: sent even if nobody guessed
	results := &RoundResult{
		RoundNb:           roundNb,
		Base:              LocationInfo{Lat: baseLat, Long: baseLong, CC: baseCC},
		PlayerRoundResult: make([]PlayerRoundResult, 0, len(activePlayers)),
	}
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
				Guess:           LocationInfo{Lat: guessLat, Long: guessLong, CC: g.CC},
				Streak:          g.Streak,
				TimeLeftSeconds: int64(g.TimeLeft.Seconds()),
			}
		}
		results.PlayerRoundResult = append(results.PlayerRoundResult, result)
	}
	return results, nil
}

// Builds the results of every round of the game from round n to round 1
func (h *GuessWsHandler) BuildGameResults(
	ctx context.Context,
	gameID, upToRoundNb int64,
	activePlayers []*games.Lobby,
) ([]RoundResult, error) {
	rounds := make([]RoundResult, 0, upToRoundNb)
	for roundNb := upToRoundNb; roundNb >= 1; roundNb-- {
		roundID, err := h.ResolveRoundID(ctx, gameID, roundNb)
		if err != nil {
			return nil, err
		}
		roundGuesses, err := h.GetRoundGuesses(ctx, roundID)
		if err != nil {
			return nil, err
		}
		result, err := h.BuildRoundResults(ctx, gameID, roundID, roundNb, activePlayers, roundGuesses)
		if err != nil {
			return nil, err
		}
		rounds = append(rounds, *result)
	}
	return rounds, nil
}

// Tells whether userID has already sent his guess for the round roundNb of gameID
// (used when a player reconnects during a round)
func (h *GuessWsHandler) HasGuessed(ctx context.Context, gameID, roundNb, userID int64) (bool, error) {
	roundID, err := h.roundRepo.GetIDByGameIDAndRoundNb(ctx, gameID, roundNb)
	if err != nil {
		return false, err
	}
	roundGuesses, err := h.guessRepo.GetByRoundID(ctx, roundID)
	if err != nil {
		return false, err
	}
	for _, g := range roundGuesses {
		if g.UserID == userID {
			return true, nil
		}
	}
	return false, nil
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

	// never create a round beyond game.round_total (CheckLastRound returns an error in that case)
	if _, err := h.roundRepo.CheckLastRound(ctx, gameID, roundNb); err != nil {
		return nil, err
	}

	//get the locationID for the round, to create the round
	locationID, err := h.roundRepo.GetLocationID(ctx, gameID, roundNb)
	if err != nil {
		return nil, err
	}

	return h.roundRepo.Create(ctx, &rounds.Round{
		GameId:     gameID,
		RoundNb:    roundNb,
		LocationId: locationID,
		StartedAt:  time.Now(),
	})
}