package handlers

import (
	"context"
	"net/http"
	"errors"
	"log"
	"strconv"
	"github.com/gin-gonic/gin"
	"backend/internal/tools/haversine"
	"backend/internal/models/guesses"
	"backend/internal/repositories"
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
	GetScoreByID(ctx context.Context, guess_id int64) (int64, error)
	GetLocationIDByRoundID(ctx context.Context, round_id int64) (int64, error)
	// GetUserIDByID(ctx context.Context, guess_id int64) (int64, error)
	// GetRoundIDByID(ctx context.Context, guess_id int64) (int64, error)
	// GetByRoundID(ctx context.Context, round_id int64) ([]*guesses.Guess, error)
	// GetByUserID(ctx context.Context, user_id int64) ([]*guesses.Guess, error)

	// PatchByID(ctx context.Context, guess_id int64, patch *guesses.Guess) error

	DeleteByID(ctx context.Context, guess_id int64) error
	// DeleteByRoundID(ctx context.Context, round_id int64) error
	// DeleteByUserID(ctx context.Context, user_id int64) error
}

type GuessHandler struct {
	guessRepo GuessGetter;
	roundRepo RoundGetter;
	lobbyRepo *repositories.LobbyRepository;
	reverser *rgeo.Rgeo;
}

func NewGuessHandler(guessRepo GuessGetter, roundRepo RoundGetter, lobbyRepo *repositories.LobbyRepository, reverser *rgeo.Rgeo) *GuessHandler {
	return &GuessHandler{guessRepo: guessRepo, roundRepo: roundRepo, lobbyRepo: lobbyRepo, reverser: reverser}
}

// Calculates distance based on lat + long compared to actual location
func calculateDistance(location, guessCoords haversine.Coordinates) (int64) {
	m, ft := haversine.Distance(location, guessCoords)
	_ = ft // currently not used, but can be returned if needed
	//TODO add an option to choose between m and ft in the request body

	dist := int64(m) // convert to int64 for storage
	return dist
}

// Calculate score based on distance and multipliers/dividers (time, cc, streaks)
func calculateScore(createdGuess *guesses.Guess) {
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

	//multiplier time left (0 to 1, based on time max being one, and time min being zero, and time left being the difference between the time the time at the moment of validating guess and the time max)
	totalTime := int64(120000000000) // TODO:  get total duration in nanoseconds (from game table where it is stored in seconds) here: maxtime to be taken from game table !
	createdGuess.TimeLeft = 120000000000 //TODO: calculate based on StoppedAt time.Time (given when user guesses , and StartedAt time.Time stored in round table)
	createdGuess.Score *= createdGuess.TimeLeft.Nanoseconds() // not in one line because it is not a float64, and does not handle decimals
	createdGuess.Score /= totalTime

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

// TODO ajouter un baseScore vs multiplierScore pour garder le baseScore pour le calcul des XP (xp = total de tous les baseScore de toutes les games jouees)

// POST /guesses
func (h *GuessHandler) AddGuess(c *gin.Context) {
	ctx := c.Request.Context()

	var guess guesses.Guess
	if err := c.ShouldBindJSON(&guess); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if guess.Lat == nil || guess.Long == nil {
    	c.JSON(http.StatusBadRequest, gin.H{"error": "latitude and longitude must be provided"})
    	return
	}
	
	if guess.RoundNb == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "round_nb must be provided"})
		return
	} // round_nb permet de recuperer le round_id en coordonnant avec le user_id (et checker si la game est en cours, et si le user a deja guess pour ce round, et si le round est bien celui de la game en cours)

	if guess.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id must be provided"})
		return // a coordonner avec le round_nb
	}

	// check if the user is in an active lobby (to get the game_id and check if the round_nb is valid for this game)
	lobby, err := h.lobbyRepo.ActualLobby(ctx, guess.UserID)
	if err != nil {
		log.Printf("Error getting active lobby for user %d: %v", guess.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get active lobby"})
		return
	}
	if lobby == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not in an active lobby"})
		return
	}

	// check if the round_nb exists in table rounds for this specific game id (if it is the next round to be played)
	guess.RoundID, err = h.roundRepo.GetIDByGameIDAndRoundNb(ctx, lobby.GameId, guess.RoundNb)
	if err != nil {
		log.Printf("Error getting round id for game %d and round_nb %d: %v", lobby.GameId, guess.RoundNb, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get round id"})
		return
	}
	
	// TODO check if game is not finished (finished value in the game table should be false, OR check if round_nb is less than round_total in the game table)
	// TODO check if the player has not left the game
	//TODO check if round is the actual round for this game (if round_nb is equal to the last round created for this game_id)

	// get locationID using roundID (to get base coordinates and base country code)
	locationID, err := h.guessRepo.GetLocationIDByRoundID(ctx, guess.RoundID)
	if err != nil {
		log.Printf("Error getting location ID for user %d at round %d: %v", guess.UserID, guess.RoundID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get location ID for round"})
		return
	}
	// get base coordinates using locationID
	log.Printf("Location ID for round %d: %d", guess.RoundID, locationID) // TODO: remove this log after testing
	baseLat, baseLong, err := h.roundRepo.GetCoordinatesByLocationID(ctx, locationID)
	if err != nil {
		log.Printf("Error getting base coordinates for round %d: %v", guess.RoundID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get base coordinates for round"})
		return
	}
	log.Printf("Base coordinates for round %d: (%f, %f)", guess.RoundID, baseLat, baseLong) // TODO: remove this log after testing
	log.Printf("Guess coordinates: (%f, %f)", *guess.Lat, *guess.Long) // TODO: remove this log after testing
	//get base cc using locationID
	baseCC, err := h.roundRepo.GetCCByLocationID(ctx, locationID)
	if err != nil {
		log.Printf("Error getting country code for location %d: %v", locationID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get country code for location"})
		return
	}
	log.Printf("Country code for location %d: %s", locationID, baseCC) // TODO: remove this log after testing

	//get cc using guess coordinates and package georeverse
	loc, err := h.reverser.ReverseGeocode([]float64{*guess.Long, *guess.Lat})
	if err != nil {
		log.Printf("Error reverse geocoding guess coordinates (%f, %f): %v", *guess.Lat, *guess.Long, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not reverse geocode guess coordinates"})
		return
	}
	guess.CC = loc.CountryCode2
	log.Printf("Country code for this guess: %s", guess.CC) // TODO: remove this log after testing
	if guess.CC == baseCC {
		guess.CountryFound = true
	} 

	// Calculate distance and score before persisting
	guess.Distance = calculateDistance(haversine.Coordinates{Lat: baseLat, Lon: baseLong}, haversine.Coordinates{Lat: *guess.Lat, Lon: *guess.Long})
	calculateScore(&guess)

	createdGuess, err := h.guessRepo.Create(ctx, &guess)
	if err != nil {
		log.Printf("Error creating guess: %v", err)
		// to check if it is a unique constraint violation (means user already guessed for this round)
		if errors.Is(err, repositories.ErrGuessNotFound) || err.Error() == "pq: duplicate key value violates unique constraint" {
			c.JSON(http.StatusConflict, gin.H{"error": "user has already submitted a guess for this round"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create guess"})
		return
	}

	c.JSON(http.StatusCreated, createdGuess)
}

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
