package handlers

import (
	// "fmt"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"backend/internal/models/games"
	"backend/internal/repositories"
)

type GameHandler struct {
	gameRepo *repositories.GameRepository
} 

func NewGameHandler(gameRepo *repositories.GameRepository) *GameHandler {
	return &GameHandler{gameRepo: gameRepo}
}

func (h *GameHandler) CreateGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		timer, err := strconv.ParseInt(c.PostForm("timer"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "the creation of the game failed, try again"})
			return ;
		}
		round_nb, err := strconv.ParseInt(c.PostForm("round_nb"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "the creation of the game failed, try again"})
			return ;
		}
		player_nb, err := strconv.ParseInt(c.PostForm("player_nb"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "the creation of the game failed, try again"})
			return ;
		}
		if timer > 2 {
			timer = 2
		}
		if round_nb > 2 {
			round_nb = 2
		}
		if player_nb > 2 {
			player_nb = 2
		}
		gameObj := &games.Game{Title: "default",
			Mode: c.PostForm("mode"),
			Timer: timer,
			RoundNb: round_nb,
			Finished: false,
			PlayerNb: player_nb}
		game, err := h.gameRepo.CreateGame(c.Request.Context(), gameObj)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "the creation of the game failed, try again"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "the game is launching", "game_id": game.ID})
	}
}

func (h *GameHandler) GetGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		game, err := h.gameRepo.GetGame(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "the game doesn't exist"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "game": game})
	}
}

func (h *GameHandler) GetTitle() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		title, err := h.gameRepo.GetTitle(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "title": title})
	}
}

func (h *GameHandler) GetMode() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		mode, err := h.gameRepo.GetMode(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "mode": mode})
	}
}

func (h *GameHandler) GetTimer() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		timer, err := h.gameRepo.GetTimer(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "timer": timer})
	}
}

func (h *GameHandler) GetRoundNb() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		roundNb, err := h.gameRepo.GetRoundNb(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "roundNb": roundNb})
	}
}

func (h *GameHandler) GetFinished() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		finished, err := h.gameRepo.GetFinished(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "finished": finished})
	}
}

func (h *GameHandler) GetPlayerNb() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		playerNb, err := h.gameRepo.GetPlayerNb(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "playerNb": playerNb})
	}
}

func (h *GameHandler) GetCreatedAt() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		createdAt, err := h.gameRepo.GetCreatedAt(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "createdAt": createdAt})
	}
}