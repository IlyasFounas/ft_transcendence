package handlers

import (
	"fmt"
	"log"
	"strconv"
	"net/http"
	"backend/internal/models/games"
	"backend/internal/repositories"
	"backend/internal/tools/code_generation"
	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameRepo *repositories.GameRepository
	lobbyRepo *repositories.LobbyRepository
	locationRepo *repositories.LocationRepository
} 

func NewGameHandler(gameRepo *repositories.GameRepository,
	lobbyRepo *repositories.LobbyRepository, 
	locationRepo *repositories.LocationRepository) *GameHandler {
	return &GameHandler{gameRepo: gameRepo, lobbyRepo: lobbyRepo, locationRepo: locationRepo}
}

func initGameObj(c *gin.Context, locationRepo *repositories.LocationRepository) (*games.Game, error) {
	timer, err := strconv.ParseInt(c.PostForm("timer"), 10, 64)
	if err != nil {
		return nil, err;
	}
	round_total, err := strconv.ParseInt(c.PostForm("round_total"), 10, 64)
	if err != nil {
		return nil, err;
	}
	if timer > 120 {
		timer = 120
	} else if (timer < 5) {
		timer = 5
	}
	if round_total > 5 {
		round_total = 5
	}
	if round_total < 1 {
		round_total = 1
	}

	var locationIDs []int64
	// check if there are enough unique locations available for the requested number of rounds
	err = locationRepo.CheckEnoughLocations(c, round_total)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, err
	}
	// get location list based on total nb of rounds
	locationIDs, err = locationRepo.GetRandomLocations(c.Request.Context(), round_total)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, err;
	}

	// c.PostForm("mode")
	gameObj := &games.Game{Title: "default",
		Mode: "Default",
		Difficulty: "Default",
		Timer: timer,
		RoundTotal: round_total,
		FinishedAt: nil,
		PlayerNb: 10,
		LocationIDs: locationIDs}
	return gameObj, nil
}

func (h *GameHandler) CreateGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		gameObj, err := initGameObj(c, h.locationRepo)
		if err != nil {
			log.Printf("Error initializing game object: %v", err)
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "the game initiation failed"})
			return ;
		}

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

func (h *GameHandler) CreateLobby() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		actualLobby, err := h.lobbyRepo.ActualLobby(c.Request.Context(), 
			userId.(int64))
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "Something went wrong"})
			return
		}
		if actualLobby != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "You're already in a lobby"})
			return
		}
		lobbyObj := &games.Lobby{UserId: userId.(int64),
								 IsCreator: true}
		for i := 0; i < 50; i++ {
			joinCode, err := code_generation.RandGenerate(4)
			if err != nil {
				c.JSON(http.StatusUnauthorized,
					gin.H{"message": "failed"})
				return
			}
			lobbyObj.JoinCode = &joinCode
			lobby, err := h.lobbyRepo.CreateLobby(c.Request.Context(), lobbyObj)
			if err == nil {
				gameObj, err := initGameObj(c, h.locationRepo)
				if err != nil {
					log.Printf("error : %s\n", err)
					h.lobbyRepo.RemoveLobby(c.Request.Context(), lobby)
					c.JSON(http.StatusUnauthorized,
						gin.H{"message": "the game initiation failed"})
					return
				}
				game, err := h.gameRepo.CreateGame(c.Request.Context(), gameObj)
				if err != nil {
					h.lobbyRepo.RemoveLobby(c.Request.Context(), lobby)
					c.JSON(http.StatusUnauthorized,
						gin.H{"message": "the creation of the game failed, try again"})
					return
				}
				err = h.lobbyRepo.SetGameId(c.Request.Context(), lobby.ID, game.ID)
				if err != nil {
					h.lobbyRepo.RemoveLobby(c.Request.Context(), lobby)
					c.JSON(http.StatusUnauthorized,
						gin.H{"message": "the lobby can not be updated, try again"})
					return
				}
				c.JSON(http.StatusOK,
					gin.H{"message": "success", "lobby_id": lobby.ID})
				return
			}
			fmt.Println("The join code already exist, currently retrying...")				
		}
		c.JSON(http.StatusUnauthorized,
			gin.H{"message": "the creation of the lobby failed, try again"})
	}
}

// Getters
func (h *GameHandler) Game() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		game, err := h.gameRepo.Game(c.Request.Context(), id)
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
		title, err := h.gameRepo.Title(c.Request.Context(), id)
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
		mode, err := h.gameRepo.Mode(c.Request.Context(), id)
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
		difficulty, err := h.gameRepo.Difficulty(c.Request.Context(), id)
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
		timer, err := h.gameRepo.Timer(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success", "timer": timer})
	}
}

func (h *GameHandler) RoundTotal() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		roundTotal, err := h.gameRepo.RoundTotal(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success", "roundTotal": roundTotal})
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
		finished, err := h.gameRepo.FinishedAt(c.Request.Context(), id)
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
		playerNb, err := h.gameRepo.PlayerNb(c.Request.Context(), id)
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
		createdAt, err := h.gameRepo.CreatedAt(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success", "createdAt": createdAt})
	}
}

// Setters
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
		err = h.gameRepo.SetTitle(c.Request.Context(), id, title)
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
		err = h.gameRepo.SetMode(c.Request.Context(), id, mode)
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
		err = h.gameRepo.SetDifficulty(c.Request.Context(), id, difficulty)
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
				gin.H{"success": false})
			return ;
		}
		err = h.gameRepo.SetFinished(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"success": false})
			return ;
		}
		// lobbies, err := h.lobbyRepo.AllPlayers(c.Request.Context(), id)
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized,
		// 		gin.H{"success": false})
		// 	return ;
		// }
		// for _, lobby := range lobbies {
		// 	if lobby.JoinCode != nil {
		// 		err = h.lobbyRepo.SetJoinCode(c.Request.Context(), lobby.ID, nil)
		// 		if err != nil {
		// 			c.JSON(http.StatusUnauthorized, 
		// 				gin.H{"success": false})
		// 			return
		// 		}
		// 	}
		// 	if lobby.LeftAt == nil {
		// 		err = h.lobbyRepo.SetLeftAt(c.Request.Context(), lobby.ID)
		// 		if err != nil {
		// 			c.JSON(http.StatusUnauthorized, 
		// 				gin.H{"success": false})
		// 			return
		// 		}
		// 	}
		// }
		c.JSON(http.StatusOK,
			gin.H{"success": true})
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
		err = h.gameRepo.SetPlayerNb(c.Request.Context(), id, player_nb)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return ;
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success"})
	}
}