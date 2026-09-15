package ws_handlers

import (
	"backend/internal/handlers"
	"backend/internal/models/guesses"
	"backend/internal/repositories"
	"backend/internal/tools/haversine"
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/sams96/rgeo"
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

type RoundInfos struct {
	Distance     int64         `json:"distance"`
	Score        int64         `json:"score"`
	BaseScore    int64         `json:"base_score"`
	CountryFound bool          `json:"country_found"`
	BaseCC       string        `json:"base_cc"`
	GuessedCC    string        `json:"guessed_cc"`
	Streak       int64         `json:"streak"`
	TimeLeft     time.Duration `json:"time_left"`
}

// résultat d'un joueur pour un round donné, avec son score total à jour dans la partie.
type PlayerRoundResult struct {
	UserID     int64      `json:"user_id"`
	Username   string     `json:"username"`
	RoundNb    int64      `json:"round_nb"`
	RoundInfos RoundInfos `json:"round_infos"`
	TotalScore int64      `json:"total_score"`
}

// valide et persiste le guess envoyé par un client WebSocket.
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
	if guess.RoundNb == 0 {
		log.Println("error : round_nb must be provided")
		return nil, false
	}
	if guess.UserID == 0 {
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
	if finished {
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
	if guess.CC == baseCC {
		guess.CountryFound = true
	}

	guess.Distance = handlers.CalculateDistance(
		haversine.Coordinates{Lat: baseLat, Lon: baseLong},
		haversine.Coordinates{Lat: *guess.Lat, Lon: *guess.Long})
	handlers.CalculateScore(&guess)

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

// renvoie tous les guess déjà soumis pour roundID
func (h *GuessWsHandler) GetRoundGuesses(ctx context.Context, roundID int64) ([]*guesses.Guess, error) {
	return h.guessRepo.GetByRoundID(ctx, roundID)
}

// transforme les guess bruts d'un round en payload de diffusion :
// un résultat détaillé par joueur, avec son score total à jour dans la partie
func (h *GuessWsHandler) BuildRoundResults(ctx context.Context, gameID int64, roundGuesses []*guesses.Guess) ([]PlayerRoundResult, error) {
	if len(roundGuesses) == 0 {
		return []PlayerRoundResult{}, nil
	}

	locationID, err := h.guessRepo.GetLocationIDByRoundID(ctx, roundGuesses[0].RoundID)
	if err != nil {
		return nil, err
	}
	baseCC, err := h.roundRepo.GetCCByLocationID(ctx, locationID)
	if err != nil {
		return nil, err
	}

	results := make([]PlayerRoundResult, 0, len(roundGuesses))
	for _, g := range roundGuesses {
		user, err := h.userRepo.User(ctx, g.UserID)
		if err != nil {
			log.Printf("error : %s\n", err)
			continue
		}
		totalScore, err := h.guessRepo.GetTotalScoreByGameAndUserID(ctx, gameID, g.UserID)
		if err != nil {
			log.Printf("error : %s\n", err)
			continue
		}
		results = append(results, PlayerRoundResult{
			UserID:   g.UserID,
			Username: user.Username,
			RoundNb:  g.RoundNb,
			RoundInfos: RoundInfos{
				Distance:     g.Distance,
				Score:        g.Score,
				BaseScore:    g.BaseScore,
				CountryFound: g.CountryFound,
				BaseCC:       baseCC,
				GuessedCC:    g.CC,
				Streak:       g.Streak,
				TimeLeft:     g.TimeLeft,
			},
			TotalScore: totalScore,
		})
	}
	return results, nil
}