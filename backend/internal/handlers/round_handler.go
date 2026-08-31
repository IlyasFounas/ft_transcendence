package handlers

import (
	"context"
	"errors"
	"log"
	"github.com/gin-gonic/gin"
	"backend/internal/models/rounds"
	"backend/internal/repositories"
)

// RoundGetter — methods expected from the repository by the handler
type RoundGetter interface {
	Create(ctx context.Context, round *rounds.Round) (*rounds.Round, error)
	
	GetByID(ctx context.Context, round_id int64) (*rounds.Round, error)

	// PatchByID(ctx context.Context, round_id int64, patch *rounds.Round) error

	DeleteByID(ctx context.Context, round_id int64) error
}

type RoundHandler struct {
	repo RoundGetter
}

func NewRoundHandler(repo RoundGetter) *RoundHandler {
	return &RoundHandler{repo: repo}
}

// POST /rounds
func (h *RoundHandler) NewRound(c *gin.Context) {
	ctx := c.Request.Context()

	var round rounds.Round
	if err := c.ShouldBindJSON(&round); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// get the game from the context
	game, exists := c.Get("game")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game not found in context"})
		return
	}

	// get the user_id from the context (from the token)
	user_id, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
		return
	}

	// set the game_id and user_id in the round struct
	round.GameId = game.(int64)
	round.UserId = user_id.(int64)

	createdRound, err := h.repo.Create(ctx, &round)
	if err != nil {
		log.Printf("Error creating round: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create round"})
		return
	}

	c.JSON(http.StatusCreated, createdRound)
}
