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
	repo GuessGetter
}

func NewGuessHandler(repo GuessGetter) *GuessHandler {
	return &GuessHandler{repo: repo}
}

// POST /guesses
func (h *GuessHandler) AddGuess(c *gin.Context) {
	ctx := c.Request.Context()

	var guess guesses.Guess
	if err := c.ShouldBindJSON(&guess); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	createdGuess, err := h.repo.Create(ctx, &guess)
	if err != nil {
		log.Printf("Error creating guess: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create guess"})
		return
	}

	if createdGuess.Lat == nil || createdGuess.Long == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "latitude and longitude must be provided"})
        return
    }
	
	// Calculates distance based on lat + long compared to actual location
	base := haversine.Coordinates{Lat: 51.45, Lon: 1.15}  // TODO replace with actual location coordinates from the database using the round_id>location_id
    guessCoords := haversine.Coordinates{Lat: *createdGuess.Lat, Lon: *createdGuess.Long,}

    m, ft := haversine.Distance(base, guessCoords)
	_ = ft // currently not used, but can be returned if needed
	//TODO add an option to choose between m and ft in the request body

	dist := int64(m) // convert to int64 for storage
	createdGuess.Distance = &dist // store distance in meters

	// TODO calculate score based on distance and multipliers/dividers (time, cc, streaks)
	
	score := haversine.EarthDiameterMeters - m
	score = score / haversine.EarthDiameterMeters * 1000 //put score on base 1000 (ex: if the distance is 0, the score is 1000, if the distance is 12742176, the score is 0)
	createdGuess.Score = int64(score)

	//multiplier time left (0 to 1, based on time max being one, and time min being zero, and time left being the difference between the time the time at the moment of validating guess and the time max)
	// TODO implement multiplier time left

	createdGuess.CountryFound = true
	// multiplier correct country code
	if createdGuess.CountryFound == true {
		createdGuess.Score *= 2
	}

	// multiplier for streaks (if less than 1000 meters, then add 1 to the streak, if more than 1000 meters, then reset the streak to 0)
	if dist < 1000 && createdGuess.CountryFound == true {
		createdGuess.Streak += 1
		createdGuess.Score *= createdGuess.Streak
	} else if (dist >= 1000 || createdGuess.CountryFound == false) {
		createdGuess.Streak = 0
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

	loc, err := h.repo.GetByID(ctx, guessID)
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

	lat, err := h.repo.GetLatitudeByID(ctx, guessID)
	if err != nil {
		if errors.Is(err, repositories.ErrGuessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "guess not found"})
			return
		}
		log.Printf("Error getting latitude for guess %d: %v", guessID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get latitude for guess"})
		return
	}

	long, err := h.repo.GetLongitudeByID(ctx, guessID)
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

	lat, err := h.repo.GetLatitudeByID(ctx, guessID)
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

	long, err := h.repo.GetLongitudeByID(ctx, guessID)
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

	cc, err := h.repo.GetCountryCodeByID(ctx, guessID)
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

	dist, err := h.repo.GetDistanceByID(ctx, guessID)
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

	score, err := h.repo.GetScoreByID(ctx, guessID)
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

	err = h.repo.DeleteByID(ctx, guessID)
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
