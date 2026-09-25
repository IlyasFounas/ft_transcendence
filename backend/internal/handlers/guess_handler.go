package handlers

import (
	"backend/internal/models/guesses"
	"backend/internal/models/rounds"
	"backend/internal/repositories"
	"backend/internal/tools/haversine"
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sams96/rgeo"
)

// GuessGetter — methods expected from the repository by the handler
type GuessGetter interface {
	Create(ctx context.Context, guess *guesses.Guess) (*guesses.Guess, error)

	GetByID(ctx context.Context, guess_id int64) (*guesses.Guess, error)
	GetLatitudeByID(ctx context.Context, guess_id int64) (float64, error)
	GetLongitudeByID(ctx context.Context, guess_id int64) (float64, error)
	GetCountryCodeByID(ctx context.Context, guess_id int64) (string, error)
	GetDistanceByID(ctx context.Context, guess_id int64) (int64, error)
	GetLocationIDByRoundID(ctx context.Context, round_id int64) (int64, error)
	GetScoreByID(ctx context.Context, guess_id int64) (int64, error)
	GetScoreByRoundAndUserID(ctx context.Context, round_id int64, user_id int64) (int64, error)
	GetTotalScoreByGameAndUserID(ctx context.Context, game_id int64, user_id int64) (int64, error)
	GetXPByUserID(ctx context.Context, user_id int64) (int64, error)
	GetResultsByGameID(ctx context.Context, game_id int64) ([][2]int64, error)
	GetUserLeaderboardInfos(ctx context.Context, user_id int64) ([]repositories.Leaderboard, error)
	GetLeaderboard(ctx context.Context) ([]repositories.Leaderboard, error)
	// GetUserIDByID(ctx context.Context, guess_id int64) (int64, error)
	// GetRoundIDByID(ctx context.Context, guess_id int64) (int64, error)
	GetByRoundID(ctx context.Context, round_id int64) ([]*guesses.Guess, error)
	// GetByUserID(ctx context.Context, user_id int64) ([]*guesses.Guess, error)

	// PatchByID(ctx context.Context, guess_id int64, patch *guesses.Guess) error

	DeleteByID(ctx context.Context, guess_id int64) error
	// DeleteByRoundID(ctx context.Context, round_id int64) error
	// DeleteByUserID(ctx context.Context, user_id int64) error
}

// RoundGetter — methods expected from the repository by the handler
type RoundGetter interface {
	Create(ctx context.Context, round *rounds.Round) (*rounds.Round, error)
	GetGame(ctx context.Context, game_id int64) error
	GetByID(ctx context.Context, round_id int64) (*rounds.Round, error)
	GetNextRoundNb(ctx context.Context, game_id int64) (int64, error)
	CheckLastRound(ctx context.Context, game_id int64, round_nb int64) (bool, error)
	GetLocationID(ctx context.Context, game_id int64, round_nb int64) (int64, error)
	GetCoordinatesByLocationID(ctx context.Context, location_id int64) (float64, float64, error)
	GetCCByLocationID(ctx context.Context, location_id int64) (string, error)
	GetIDByGameIDAndRoundNb(ctx context.Context, game_id int64, round_nb int64) (int64, error)
	// DeleteByID(ctx context.Context, round_id int64) error
}

type GuessHandler struct {
	guessRepo	GuessGetter
	roundRepo	RoundGetter
	lobbyRepo	*repositories.LobbyRepository
	gameRepo	*repositories.GameRepository
	reverser	*rgeo.Rgeo
}

func NewGuessHandler(guessRepo GuessGetter, roundRepo RoundGetter, lobbyRepo *repositories.LobbyRepository,
	gameRepo *repositories.GameRepository, reverser *rgeo.Rgeo) *GuessHandler {
	return &GuessHandler{
		guessRepo:	guessRepo,
		roundRepo:	roundRepo,
		lobbyRepo:	lobbyRepo,
		gameRepo:	gameRepo,
		reverser:	reverser}
}

// Calculates distance based on lat + long compared to actual location
func CalculateDistance(location, guessCoords haversine.Coordinates) int64 {
	m, ft := haversine.Distance(location, guessCoords)
	_ = ft // currently not used, but can be returned if needed
	//TODO add an option to choose between m and ft in the request body

	dist := int64(m) // convert to int64 for storage
	return dist
}

// Calculate score based on distance and multipliers/dividers (time, cc, streaks)
func CalculateScore(createdGuess *guesses.Guess, totalTime time.Duration) {
	if createdGuess.Distance < 0 {
		log.Printf("Error: negative distance %d for guess", createdGuess.Distance)
		createdGuess.Score = 0
		return
	}

	if createdGuess.Distance > haversine.EarthDiameterMeters {
		log.Printf("Error: distance %d exceeds Earth's diameter for guess", createdGuess.Distance)
		createdGuess.Score = 0
		return
	}

	distance := haversine.EarthDiameterMeters - createdGuess.Distance
	createdGuess.Score = distance * 1000 / haversine.EarthDiameterMeters //put score on base 1000 (ex: if the distance is 0, the score is 1000, if the distance is 12742176, the score is 0)
	createdGuess.BaseScore = createdGuess.Score // store the base score before multipliers for XP calculation

	// multiplier time left (0 to 1, based on time max being one, and time min being zero), and time left being the difference between the time the time at the moment of validating guess and the time max)
	log.Printf("Score before time left multiplier: %d (TimeLeft: %v, totalTime: %v)\n", createdGuess.Score, createdGuess.TimeLeft, totalTime)
	createdGuess.Score *= createdGuess.TimeLeft.Nanoseconds() // not in one line because it is not a float64, and does not handle decimals
	log.Printf("Score after time left multiplier: %d (TimeLeft: %v, totalTime: %v)\n", createdGuess.Score, createdGuess.TimeLeft, totalTime)
	if totalTime > 0 {
		createdGuess.Score /= totalTime.Nanoseconds()
	}

	// multiplier correct country code
	if createdGuess.CountryFound == true {
		createdGuess.Score *= 2
	}

	// multiplier for streaks (if less than 1000 meters and country found, then add 1 to the streak, if more than 1000 meters, if at least one is true, then keep streak, if none then reset the streak to 0)
	if createdGuess.Distance < 1000 && createdGuess.CountryFound == true {
		createdGuess.Streak += 1
		createdGuess.Score *= createdGuess.Streak
	} else if createdGuess.Distance >= 1000 || createdGuess.CountryFound == false {
		createdGuess.Streak = 0
	}
}

// // POST /guesses
// func (h *GuessHandler) AddGuess(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	var guess guesses.Guess
// 	if err := c.ShouldBindJSON(&guess); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
// 		return
// 	}

// 	if guess.Lat == nil || guess.Long == nil {
//     	c.JSON(http.StatusBadRequest, gin.H{"error": "latitude and longitude must be provided"})
//     	return
// 	}

// 	if guess.RoundNb == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "round_nb must be provided"})
// 		return
// 	} // round_nb permet de recuperer le round_id en coordonnant avec le user_id (et checker si la game est en cours, et si le user a deja guess pour ce round, et si le round est bien celui de la game en cours)

// 	if guess.UserID == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id must be provided"})
// 		return // a coordonner avec le round_nb
// 	}

// 	// check if the user is in an active lobby (to get the game_id and check if the round_nb is valid for this game)
// 	lobby, err := h.lobbyRepo.ActualLobby(ctx, guess.UserID)
// 	if err != nil {
// 		log.Printf("Error getting active lobby for user %d: %v", guess.UserID, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get active lobby"})
// 		return
// 	}
// 	if lobby == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not in an active lobby"})
// 		return
// 	}

// 	// check if the round_nb exists in table rounds for this specific game id (if it is the next round to be played)
// 	guess.RoundID, err = h.roundRepo.GetIDByGameIDAndRoundNb(ctx, lobby.GameId, guess.RoundNb)
// 	if err != nil {
// 		log.Printf("Error getting round id for game %d and round_nb %d: %v", lobby.GameId, guess.RoundNb, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get round id"})
// 		return
// 	}

// 	// check if the game is finished (if true, then the user cannot guess anymore) OR check if round_nb is less than round_total in the game table
// 	finished, err := h.gameRepo.Finished(ctx, lobby.GameId)
// 	if err != nil {
// 		log.Printf("Error checking if game %d is finished: %v", lobby.GameId, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not check game status"})
// 		return
// 	}
// 	if finished {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "this game is already finished"})
// 		return
// 	}

// 	// TODO check if the player has not left the lobby

// 	// round_nb should be the same as the current round (last round created for the game), else the user is trying to guess for a round that is not the current one (or has already been played)
// 	nextRoundNb, err := h.roundRepo.GetNextRoundNb(ctx, lobby.GameId)
// 	if err != nil {
// 		log.Printf("Error getting current round number for game %d: %v", lobby.GameId, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not verify round number"})
// 		return
// 	}
// 	if guess.RoundNb != nextRoundNb - 1 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "round_nb does not match the current round for this game"})
// 		return
// 	}

// 	// get locationID using roundID (to get base coordinates and base country code)
// 	locationID, err := h.guessRepo.GetLocationIDByRoundID(ctx, guess.RoundID)
// 	if err != nil {
// 		log.Printf("Error getting location ID for user %d at round %d: %v", guess.UserID, guess.RoundID, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get location ID for round"})
// 		return
// 	}
// 	// get base coordinates using locationID
// 	log.Printf("Location ID for round %d: %d", guess.RoundID, locationID) // TODO: remove this log after testing
// 	baseLat, baseLong, err := h.roundRepo.GetCoordinatesByLocationID(ctx, locationID)
// 	if err != nil {
// 		log.Printf("Error getting base coordinates for round %d: %v", guess.RoundID, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get base coordinates for round"})
// 		return
// 	}
// 	log.Printf("Base coordinates for round %d: (%f, %f)", guess.RoundID, baseLat, baseLong) // TODO: remove this log after testing
// 	log.Printf("Guess coordinates: (%f, %f)", *guess.Lat, *guess.Long) // TODO: remove this log after testing
// 	//get base cc using locationID
// 	baseCC, err := h.roundRepo.GetCCByLocationID(ctx, locationID)
// 	if err != nil {
// 		log.Printf("Error getting country code for location %d: %v", locationID, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get country code for location"})
// 		return
// 	}
// 	log.Printf("Country code for location %d: %s", locationID, baseCC) // TODO: remove this log after testing

// 	//get cc using guess coordinates and package georeverse
// 	loc, err := h.reverser.ReverseGeocode([]float64{*guess.Long, *guess.Lat})
// 	if err != nil {
// 		log.Printf("Error reverse geocoding guess coordinates (%f, %f): %v", *guess.Lat, *guess.Long, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not reverse geocode guess coordinates"})
// 		return
// 	}
// 	guess.CC = loc.CountryCode2
// 	log.Printf("Country code for this guess: %s", guess.CC) // TODO: remove this log after testing
// 	if guess.CC == baseCC {
// 		guess.CountryFound = true
// 	}

// 	// Calculate distance and score before persisting
// 	guess.Distance = CalculateDistance(haversine.Coordinates{Lat: baseLat, Lon: baseLong}, haversine.Coordinates{Lat: *guess.Lat, Lon: *guess.Long})
// 	CalculateScore(&guess)

// 	createdGuess, err := h.guessRepo.Create(ctx, &guess)
// 	if err != nil {
// 		log.Printf("Error creating guess: %v", err)
// 		// to check if it is a unique constraint violation (means user already guessed for this round)
// 		if errors.Is(err, repositories.ErrGuessNotFound) || err.Error() == "pq: duplicate key value violates unique constraint" {
// 			c.JSON(http.StatusConflict, gin.H{"error": "user has already submitted a guess for this round"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create guess"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, createdGuess)
// }

// GET /guesses/:id
func (h *GuessHandler) GetGuessByID(c *gin.Context) {
	ctx := c.Request.Context()

	guessID, err := strconv.ParseInt(c.Param("guess_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guess id"})
		return
	}

	loc, err := h.guessRepo.GetByID(ctx, guessID)
	if err != nil {
		if errors.Is(err, repositories.ErrGuessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "guess not found"})
			return
		}
		log.Printf("Error getting guess %d: %v", guessID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get guess"})
		return
	}

	c.JSON(http.StatusOK, loc)
}

// GET /guesses/:id/coordinates
func (h *GuessHandler) GetGuessCoordinates(c *gin.Context) {
	ctx := c.Request.Context()

	guessID, err := strconv.ParseInt(c.Param("guess_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guess id"})
		return
	}

	lat, err := h.guessRepo.GetLatitudeByID(ctx, guessID)
	if err != nil {
		if errors.Is(err, repositories.ErrGuessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "guess not found"})
			return
		}
		log.Printf("Error getting latitude for guess %d: %v", guessID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get latitude for guess"})
		return
	}

	long, err := h.guessRepo.GetLongitudeByID(ctx, guessID)
	if err != nil {
		if errors.Is(err, repositories.ErrGuessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "guess not found"})
			return
		}
		log.Printf("Error getting longitude for guess %d: %v", guessID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get longitude for guess"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"latitude": lat, "longitude": long})
}

// GET /guesses/:id/latitude
func (h *GuessHandler) GetGuessLatitude(c *gin.Context) {
	ctx := c.Request.Context()

	guessID, err := strconv.ParseInt(c.Param("guess_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guess id"})
		return
	}

	lat, err := h.guessRepo.GetLatitudeByID(ctx, guessID)
	if err != nil {
		if errors.Is(err, repositories.ErrGuessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "guess not found"})
			return
		}
		log.Printf("Error getting latitude for guess %d: %v", guessID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get latitude for guess"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"latitude": lat})
}

// GET /guesses/:id/longitude
func (h *GuessHandler) GetGuessLongitude(c *gin.Context) {
	ctx := c.Request.Context()

	guessID, err := strconv.ParseInt(c.Param("guess_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guess id"})
		return
	}

	long, err := h.guessRepo.GetLongitudeByID(ctx, guessID)
	if err != nil {
		if errors.Is(err, repositories.ErrGuessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "guess not found"})
			return
		}
		log.Printf("Error getting longitude for guess %d: %v", guessID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get longitude for guess"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"longitude": long})
}

// GET /guesses/:id/country
func (h *GuessHandler) GetGuessCountryCode(c *gin.Context) {
	ctx := c.Request.Context()

	guessID, err := strconv.ParseInt(c.Param("guess_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guess id"})
		return
	}

	cc, err := h.guessRepo.GetCountryCodeByID(ctx, guessID)
	if err != nil {
		if errors.Is(err, repositories.ErrGuessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "guess not found"})
			return
		}
		log.Printf("Error getting country code for guess %d: %v", guessID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get country code for guess"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"country_code": cc})
}

// GET /guesses/:id/distance
func (h *GuessHandler) GetGuessDistance(c *gin.Context) {
	ctx := c.Request.Context()

	guessID, err := strconv.ParseInt(c.Param("guess_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guess id"})
		return
	}

	dist, err := h.guessRepo.GetDistanceByID(ctx, guessID)
	if err != nil {
		if errors.Is(err, repositories.ErrGuessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "guess not found"})
			return
		}
		log.Printf("Error getting distance for guess %d: %v", guessID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get distance for guess"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"distance": dist})
}

// GET /guesses/:id/score
func (h *GuessHandler) GetGuessScore(c *gin.Context) {
	ctx := c.Request.Context()

	guessID, err := strconv.ParseInt(c.Param("guess_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guess id"})
		return
	}

	score, err := h.guessRepo.GetScoreByID(ctx, guessID)
	if err != nil {
		if errors.Is(err, repositories.ErrGuessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "guess not found"})
			return
		}
		log.Printf("Error getting score for guess %d: %v", guessID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get score for guess"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"score": score})
}

// GET /guesses/score/:round_id/:user_id //by round and and user
func (h *GuessHandler) GetGuessScoreByRoundAndUser(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized,
			gin.H{"message": "The userId is empty"})
		return
	}

	ctx := c.Request.Context()

	round_id, err := strconv.ParseInt(c.Param("round_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid round id"})
		return
	}

	score, err := h.guessRepo.GetScoreByRoundAndUserID(ctx, round_id, userID.(int64))
	if err != nil {
		if errors.Is(err, repositories.ErrGuessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "guess not found for this round and user"})
			return
		}
		log.Printf("Error getting score for round %d and user %d: %v", round_id, userID.(int64), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get score for this round and user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"score": score})
}

// GET /guesses/total_score/:game_id/:user_id //by game and user
func (h *GuessHandler) GetTotalScoreByGameAndUser(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized,
			gin.H{"message": "The userId is empty"})
		return
	}

	ctx := c.Request.Context()

	game_id, err := strconv.ParseInt(c.Param("game_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	totalScore, err := h.guessRepo.GetTotalScoreByGameAndUserID(ctx, game_id, userID.(int64))
	if err != nil {
		log.Printf("Error getting total score for game %d and user %d: %v", game_id, userID.(int64), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get total score for this game and user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_score": totalScore})
}

// GET /guesses/xp/:user_id //by user
func (h *GuessHandler) GetUserXP(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized,
			gin.H{"message": "The userId is empty"})
		return
	}

	ctx := c.Request.Context()

	userXP, err := h.guessRepo.GetXPByUserID(ctx, userID.(int64))
	if err != nil {
		log.Printf("Error getting total score for user %d: %v", userID.(int64), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get total score for this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_xp": userXP})
}

// GET /guesses/results/:game_id //by game
func (h *GuessHandler) GetGameResults(c *gin.Context) {
	ctx := c.Request.Context()

	game_id, err := strconv.ParseInt(c.Param("game_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	results, err := h.guessRepo.GetResultsByGameID(ctx, game_id)
	if err != nil {
		log.Printf("Error getting results for game %d: %v", game_id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get results for this game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

// GET /leaderboard/user //using user in the context + showing 100 first players and position and results of the user
func (h *GuessHandler) GetUserLeaderboardInfos(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized,
			gin.H{"message": "The user ID is empty"})
		return
	}

	leaderboard, err := h.guessRepo.GetUserLeaderboardInfos(ctx, userID.(int64))
	if err != nil {
		log.Printf("Error getting leaderboard for user %d: %v", userID.(int64), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user leaderboard infos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
}

// GET /leaderboard //using user in the context + showing 100 first players and position and results of the user
func (h *GuessHandler) GetLeaderboard(c *gin.Context) {
	ctx := c.Request.Context()

	leaderboard, err := h.guessRepo.GetLeaderboard(ctx)
	if err != nil {
		log.Printf("Error getting leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get leaderboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
}

// // GET /guesses/round/:round_id
// func (h *GuessHandler) GetGuessesByRoundID(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	round_id, err := strconv.ParseInt(c.Param("round_id"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid round id"})
// 		return
// 	}

// 	guesses, err := h.guessRepo.GetByRoundID(ctx, round_id)
// 	if err != nil {
// 		log.Printf("Error getting guesses for round %d: %v", round_id, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get guesses for this round"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"guesses": guesses})
// }

// DELETE /guesses/:id
func (h *GuessHandler) DeleteGuessByID(c *gin.Context) {
	ctx := c.Request.Context()

	guessID, err := strconv.ParseInt(c.Param("guess_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guess id"})
		return
	}

	err = h.guessRepo.DeleteByID(ctx, guessID)
	if err != nil {
		if errors.Is(err, repositories.ErrGuessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "guess not found"})
			return
		}
		log.Printf("Error deleting guess %d: %v", guessID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete guess"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "guess deleted successfully"}) // or no content ?
}