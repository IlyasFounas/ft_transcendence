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
	// GetLocationID(ctx context.Context, game_id int64, round_nb int64) (int64, error)

	// PatchByID(ctx context.Context, round_id int64, patch *rounds.Round) error

	// DeleteByID(ctx context.Context, round_id int64) error
}

// type GameGetter interface {
// 	// GetByID(ctx context.Context, game_id int64) (*rounds.Round, error)
// }

type RoundHandler struct {
	repo RoundGetter
	// repo GameGetter
}

func NewRoundHandler(repo RoundGetter) *RoundHandler {
	return &RoundHandler{repo: repo}
}

// POST /rounds
func (h *RoundHandler) NextRound(c *gin.Context) {
	log.Println("NewRound handler called\n") //TODO to remove
	ctx := c.Request.Context()

	var round rounds.Round
	// recuperer le game_id depuis le body de la requete
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
	err := h.repo.GetGame(ctx, round.GameId)
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

	// recuperer le nombre de rounds deja cree pour ce game_id + 1 pour calculer le bon round_nb actuel (une jointure avec la table rounds est necessaire pour recuperer le nombre de rounds deja crees pour ce game_id)
	round.RoundNb, err = h.repo.GetRoundNb(ctx, round.GameId)
	if err != nil {
		log.Printf("Error counting rounds: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not count rounds"})
		return
	}

	// TODO si round number superieur au nombre total de rounds present dans game_id alors on ressort une erreur pour dire que c'est la fin de la game

	// // // get the user_id from the context (from the token)
	// // user_id, exists := c.Get("userId")
	// // if !exists {
	// // 	c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
	// // 	return
	// // }
	// // // set the user_id in the round struct
	// // round.UserId = user_id.(int64)

	// // recuperer la location_id associee a ce round (avec la game_id et le round_nb), la liste de locations est dans la table games, il faut donc faire une jointure pour recuperer la location_id associee a ce round
	// location_id, err := h.repo.GetLocationID(ctx, round.GameId, round.RoundNb)
	// if err != nil {
	// 	log.Printf("Error getting location_id: %v", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get location_id"})
	// 	return
	// }
	// round.LocationId = location_id

	// // si la game_id, round_nb et location_id sont valides, on peut creer le round dans la base de donnees
	// if (round.GameId == 0 || round.RoundNb == 0 || round.LocationId == 0) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game_id, round_nb or location_id"})
	// 	return
	// }
	
	createdRound, err := h.repo.Create(ctx, &round)
	if err != nil {
		log.Printf("Error creating round: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create round"})
		fmt.Printf("5 Error create round + gameid: %v + %d\n", err, round.GameId) //TODO to remove
		return
	}

	c.JSON(http.StatusCreated, createdRound)
	fmt.Printf("6 Error round + gameid: %v + %d\n", err, round.GameId) //TODO to remove

}

