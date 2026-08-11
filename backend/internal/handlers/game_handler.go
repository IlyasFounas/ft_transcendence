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
		if timer > 120 {
			timer = 120
		} else if (timer < 10) {
			timer = 10
		}
		if round_nb > 2 {
			round_nb = 2
		}
		if player_nb > 2 {
			player_nb = 2
		}
		gameObj := &games.Game{Title: "default",
			Mode: c.PostForm("mode"),
			Difficulty: "easy",
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

func (h *GameHandler) Game() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		game, err := h.gameRepo.Game(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "the game doesn't exist"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "game": game})
	}
}

func (h *GameHandler) Title() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		title, err := h.gameRepo.Title(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "title": title})
	}
}

func (h *GameHandler) Mode() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		mode, err := h.gameRepo.Mode(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "mode": mode})
	}
}

func (h *GameHandler) Difficulty() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		difficulty, err := h.gameRepo.Difficulty(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "difficulty": difficulty})
	}
}

func (h *GameHandler) Timer() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		timer, err := h.gameRepo.Timer(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "timer": timer})
	}
}

func (h *GameHandler) RoundNb() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		roundNb, err := h.gameRepo.RoundNb(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "roundNb": roundNb})
	}
}

func (h *GameHandler) Finished() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		finished, err := h.gameRepo.Finished(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "finished": finished})
	}
}

func (h *GameHandler) PlayerNb() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		playerNb, err := h.gameRepo.PlayerNb(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "playerNb": playerNb})
	}
}

func (h *GameHandler) CreatedAt() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		createdAt, err := h.gameRepo.CreatedAt(c, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success", "createdAt": createdAt})
	}
}

func (h *GameHandler) SetTitle() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		title := c.PostForm("title")
		if (len(title) == 0) {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The title is empty"})
			return ;
		}
		err = h.gameRepo.SetTitle(c, id, title)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success"})
	}
}

func (h *GameHandler) SetMode() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		mode := c.PostForm("mode")
		if (len(mode) == 0) {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The mode is empty"})
			return ;
		}
		err = h.gameRepo.SetMode(c, id, mode)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success"})
	}
}

func (h *GameHandler) SetDifficulty() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		difficulty := c.PostForm("difficulty")
		if (len(difficulty) == 0) {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The difficulty is empty"})
			return ;
		}
		err = h.gameRepo.SetDifficulty(c, id, difficulty)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success"})
	}
}

func (h *GameHandler) SetTimer() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		timer, err := strconv.ParseInt(c.PostForm("timer"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		if timer > 120 {
			timer = 120
		} else if (timer < 10) {
			timer = 10
		}
		err = h.gameRepo.SetTimer(c, id, timer)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success"})
	}
}

func (h *GameHandler) SetRoundNb() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		round_nb, err := strconv.ParseInt(c.PostForm("round_nb"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}

		if round_nb > 5 {
			round_nb = 5
		} else if (round_nb < 1) {
			round_nb = 1
		}
		err = h.gameRepo.SetRoundNb(c, id, round_nb)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success"})
	}
}

func (h *GameHandler) SetFinished() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		finished, err := strconv.ParseBool(c.PostForm("finished"))
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		err = h.gameRepo.SetFinished(c, id, finished)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success"})
	}
}

func (h *GameHandler) SetPlayerNb() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		player_nb, err := strconv.ParseInt(c.PostForm("player_nb"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		if player_nb > 10 {
			player_nb = 10
		} else if (player_nb < 1) {
			player_nb = 1
		}
		err = h.gameRepo.SetPlayerNb(c, id, player_nb)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
        c.JSON(http.StatusOK,
			gin.H{"message": "success"})
	}
}