package handlers

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/internal/models/users"
	"backend/internal/repositories"
)

type StatisticsHandler struct {
	statsRepo *repositories.StatisticsRepository
}

func NewStatisticsHandler(statsRepo *repositories.StatisticsRepository) *StatisticsHandler {
	return &StatisticsHandler{statsRepo: statsRepo}
}

type statisticsResponse struct {
	FinishedGames int64           `json:"finished_games"`
	Accuracy      float64         `json:"accuracy"`      // percentage, rounded to 2 decimals
	BestScore     int64           `json:"best_score"`
	ClosestGuess  *users.Distance `json:"closest_guess"` // struct{value, unit: "m" or "km"} OR null if no guess yet
}

// Rounds a percentage to 2 decimals for display (the average itself is stored as is)
func roundAccuracy(accuracy float64) float64 {
	return math.Round(accuracy*100) / 100
}

// Converts the closest guess (in meters) into a distance ready to display, nil if there is none
func closestGuess(meters *int64) *users.Distance {
	if meters == nil {
		return nil
	}
	distance := users.NewDistance(*meters)
	return &distance
}

// Getters
func (h *StatisticsHandler) Statistics() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
			gin.H{"message": "The userId is empty"})
			return
		}
		stats, err := h.statsRepo.Statistics(c.Request.Context(), userId.(int64))
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success", "statistics": statisticsResponse{
			FinishedGames: stats.FinishedGames,
			Accuracy:      roundAccuracy(stats.Accuracy),
			BestScore:     stats.BestScore,
			ClosestGuess:  closestGuess(stats.ClosestGuess),
		}})
	}
}

func (h *StatisticsHandler) FinishedGames() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
			gin.H{"message": "The userId is empty"})
			return
		}
		finishedGames, err := h.statsRepo.FinishedGames(c.Request.Context(), userId.(int64))
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success", "finished_games": finishedGames})
	}
}

func (h *StatisticsHandler) Accuracy() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
			gin.H{"message": "The userId is empty"})
			return
		}
		accuracy, err := h.statsRepo.Accuracy(c.Request.Context(), userId.(int64))
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success", "accuracy": roundAccuracy(accuracy)})
	}
}

func (h *StatisticsHandler) BestScore() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
			gin.H{"message": "The userId is empty"})
			return
		}
		bestScore, err := h.statsRepo.BestScore(c.Request.Context(), userId.(int64))
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success", "best_score": bestScore})
	}
}

func (h *StatisticsHandler) ClosestGuess() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
			gin.H{"message": "The userId is empty"})
			return
		}
		meters, err := h.statsRepo.ClosestGuess(c.Request.Context(), userId.(int64))
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success", "closest_guess": closestGuess(meters)})
	}
}

// Setters
func setterResponse(c *gin.Context, updated bool, err error) {
	if err != nil {
		log.Printf("error : %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}
	if !updated {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid value"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *StatisticsHandler) SetFinishedGames() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
			return
		}
		finishedGames, err := strconv.ParseInt(c.PostForm("finished_games"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid finished_games"})
			return
		}
		updated, err := h.statsRepo.SetFinishedGames(c.Request.Context(), userId, finishedGames)
		setterResponse(c, updated, err)
	}
}

func (h *StatisticsHandler) SetAccuracy() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
			return
		}
		accuracy, err := strconv.ParseFloat(c.PostForm("accuracy"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid accuracy"})
			return
		}
		updated, err := h.statsRepo.SetAccuracy(c.Request.Context(), userId, accuracy)
		setterResponse(c, updated, err)
	}
}

func (h *StatisticsHandler) SetBestScore() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
			return
		}
		bestScore, err := strconv.ParseInt(c.PostForm("best_score"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid best_score"})
			return
		}
		updated, err := h.statsRepo.SetBestScore(c.Request.Context(), userId, bestScore)
		setterResponse(c, updated, err)
	}
}

func (h *StatisticsHandler) SetClosestGuess() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
			return
		}
		meters, err := strconv.ParseInt(c.PostForm("closest_guess"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid closest_guess"})
			return
		}
		updated, err := h.statsRepo.SetClosestGuess(c.Request.Context(), userId, meters)
		setterResponse(c, updated, err)
	}
}

// AddGameStatistics adds the game gameId to the statistics of every player still in it when it ends (not the ones
// who left it or were kicked out). It updates:
//   - finished_games: +1
//   - accuracy: the accuracy of a guess is a rule of three on its distance: 0 m = 100 %, haversine.EarthDiameterMeters
//     or more = 0 % (the same reference as the score). The accuracy of the game is the average over its rounds (a round
//     without guess counts 0 %), and the statistic is the average over all his finished games.
//   - best_score: replaced if the final score of this game is higher
//   - closest_guess: replaced if one of his guesses of this game is closer
func (r *StatisticsHandler) AddGameStatistics(c *gin.Context, gameId int64) error {
	ctx := c.Request.Context()

	// get the round total of the game
	roundTotal, err := r.statsRepo.GetRoundTotal(ctx, gameId)
	if err != nil {
		return err
	}
	if roundTotal < 1 {
		log.Printf("game %d has no rounds to add to statistics", gameId)
		return nil
	}

	// get the players that stayed till the end of the game
	players, err := r.statsRepo.GetPlayersInGame(ctx, gameId)
	if err != nil {
		log.Printf("error getting players in game %d: %s", gameId, err)
		return err
	}

	// get the statistics of each player in the game
	userGameStats, err := r.statsRepo.SetGameStatisticsPerUser(ctx, gameId)
	if err != nil {
		log.Printf("error setting game statistics for game %d: %s", gameId, err)
		return err
	}
	
	err = r.statsRepo.UpdateStatistics(ctx, userGameStats, players, roundTotal)
	if err != nil {
		log.Printf("error updating user statistics after game %d: %s", gameId, err)
		return err
	}
	return nil
}