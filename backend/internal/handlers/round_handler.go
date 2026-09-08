package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"backend/internal/models/rounds"
	"backend/internal/repositories"
)

// RoundGetter — methods expected from the repository by the handler
type RoundGetter interface {
	Create(ctx context.Context, round *rounds.Round) (*rounds.Round, error)
	GetGame(ctx context.Context, game_id int64) error
	// GetByID(ctx context.Context, round_id int64) (*rounds.Round, error)
	GetRoundNb(ctx context.Context, game_id int64) (int64, error)
	CheckLastRound(ctx context.Context, game_id int64, round_nb int64) (bool, error)
	GetLocationID(ctx context.Context, game_id int64, round_nb int64) (int64, error)
	GetCoordinatesByLocationID(ctx context.Context, location_id int64) (float64, float64, error)
	GetCCByLocationID(ctx context.Context, location_id int64) (string, error)
	// DeleteByID(ctx context.Context, round_id int64) error
}

type RoundHandler struct {
	roundRepo RoundGetter
	locationRepo *repositories.LocationRepository
}

func NewRoundHandler(roundRepo RoundGetter) *RoundHandler {
	return &RoundHandler{roundRepo: roundRepo}
}

// POST /rounds
func (h *RoundHandler) NextRound(c *gin.Context) {
	log.Println("NewRound handler called") //TODO to remove
	ctx := c.Request.Context()

	var round rounds.Round
	// retrieve the game_id from the request body
	if err := c.ShouldBindJSON(&round); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		fmt.Printf("0 Error binding JSON + gameid: %v + %d\n", err, round.GameId) //TODO to remove
		return
	}
	// TODO ou depuis le context (si on a un middleware qui le met dans le context)
	// gameid := c.Query("game_id")
	// // get the game from the context
	// game, exists := c.Get("game")
	// if !exists {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "game not found in context"})
	// 	return
	// }
	// // set the game_id in the round struct
	// round.GameId = gameid.(int64)

	// check if the game_id has been provided in the request body
	if round.GameId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game_id is required"})
		fmt.Printf("2 Error binding JSON gameid: %d %d\n", round.ID, round.GameId) //TODO to remove
		return
	}

	// check if the game_id exists in the database
	err := h.roundRepo.GetGame(ctx, round.GameId)
	if err != nil {
		if errors.Is(err, repositories.ErrGameNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "game_id does not exist"})
			fmt.Printf("3 Error get game + gameid: %v + %d\n", err, round.GameId) //TODO to remove
			return
		}
		log.Printf("Error checking game_id: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not check game_id"})
		fmt.Printf("4 Error get game + gameid: %v + %d\n", err, round.GameId) //TODO to remove
		return
	}

	// calculate the actual round number using the number of rounds already created + 1
	round.RoundNb, err = h.roundRepo.GetRoundNb(ctx, round.GameId)
	if err != nil {
		log.Printf("Error counting rounds: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not count rounds"})
		return
	}

	// check if this is the last round for this game_id
	round.LastRound, err = h.roundRepo.CheckLastRound(ctx, round.GameId, round.RoundNb)
	if err != nil {
		log.Printf("Error getting total rounds: %v", err)
		if round.LastRound == true {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// // // get the user_id from the context (from the token)
	// // user_id, exists := c.Get("userId")
	// // if !exists {
	// // 	c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
	// // 	return
	// // }
	// // // set the user_id in the round struct
	// // round.UserId = user_id.(int64)

	// retrieve the location_id for this game_id and round_nb
	location_id, err := h.roundRepo.GetLocationID(ctx, round.GameId, round.RoundNb)
	if err != nil {
		log.Printf("Error getting location_id: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get location_id"})
		return
	}
	round.LocationId = location_id
	log.Printf("Location ID for game_id %d and round_nb %d: %d", round.GameId, round.RoundNb, location_id) //TODO to remove

	// // si la game_id, round_nb et location_id sont valides, on peut creer le round dans la base de donnees
	// if (round.GameId == 0 || round.RoundNb == 0 || round.LocationId == 0) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game_id, round_nb or location_id"})
	// 	return
	// }
	
	createdRound, err := h.roundRepo.Create(ctx, &round)
	if err != nil {
		log.Printf("Error creating round: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create round"})
		fmt.Printf("5 Error create round + gameid: %v + %d\n", err, round.GameId) //TODO to remove
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"game_id":     createdRound.GameId,
		"round_nb":    createdRound.RoundNb,
		"last_round":  createdRound.LastRound,
	})
}


// TODO function to retrieve the round corresponding coordinates (to launch the game) ?