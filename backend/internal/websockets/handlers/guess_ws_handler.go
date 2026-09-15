package ws_handlers

import (
	"backend/internal/handlers"
	"backend/internal/models/games"
	"backend/internal/models/guesses"
	"backend/internal/repositories"
	"backend/internal/tools/haversine"
	"context"
	"encoding/json"
	"errors"
	"log"
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
		guessRepo:	guessRepo,
		roundRepo:	roundRepo,
		lobbyRepo:	lobbyRepo,
		gameRepo:	gameRepo,
		userRepo:	userRepo,
		reverser:	reverser,
	}
}

type LocationInfo struct {
	Lat		float64	`json:"lat"`
	Long	float64	`json:"long"`
	CC		string	`json:"country_code"`
}

type RoundInfos struct {
	Distance		int64			`json:"distance"`
	Score			int64			`json:"score"`
	BaseScore		int64			`json:"base_score"`
	CountryFound	bool			`json:"country_found"`
	Base			LocationInfo	`json:"base"`
	Guess			LocationInfo	`json:"guess"`
	Streak			int64			`json:"streak"`
	TimeLeftSeconds	int64			`json:"time_left_seconds"`
}

// résultat d'un joueur pour un round donné, avec son score total à jour dans la partie
type PlayerRoundResult struct {
	UserID		int64		`json:"user_id"`
	Username	string		`json:"username"`
	RoundNb		int64		`json:"round_nb"`
	HasGuessed	bool		`json:"has_guessed"`
	RoundInfos	*RoundInfos	`json:"round_infos,omitempty"`
	TotalScore	int64		`json:"total_score"`
}

// renvoie tous les guess déjà soumis pour roundID
func (h *GuessWsHandler) GetRoundGuesses(ctx context.Context, roundID int64) ([]*guesses.Guess, error) {
	return h.guessRepo.GetByRoundID(ctx, roundID)
}

// Transforme les guess bruts d'un round en payload de diffusion :
// un résultat détaillé par joueur, avec son score total à jour dans la partie.
// Construit un résultat pour CHAQUE joueur actif de la game, pas seulement ceux qui ont deviné.
func (h *GuessWsHandler) BuildRoundResults(ctx context.Context, gameID, roundID, roundNb int64,
	activePlayers []*games.Lobby, roundGuesses []*guesses.Guess) ([]PlayerRoundResult, error) {

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