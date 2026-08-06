package handlers

import (
	"net/htpp"
	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameRepo *repositories.gameHandler
} 

func NewGameHandler(g GameHandler) {
	return &GameHandler{gameRepo: gameRepo}
}

func (h *GameHandler) CreateGame() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}