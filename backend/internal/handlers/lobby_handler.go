package handlers

import (
	"backend/internal/models/games"
	"backend/internal/repositories"
	"log"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type LobbyHandler struct {
	gameRepo *repositories.GameRepository
	guessRepo *repositories.GuessRepository
	lobbyRepo *repositories.LobbyRepository
	userRepo *repositories.UserRepository
}

func NewLobbyHandler(gameRepo *repositories.GameRepository,
	guessRepo *repositories.GuessRepository,
	lobbyRepo *repositories.LobbyRepository,
	userRepo *repositories.UserRepository) *LobbyHandler {
	return &LobbyHandler{gameRepo: gameRepo, guessRepo: guessRepo, lobbyRepo: lobbyRepo,
		userRepo: userRepo}
}

func (h *LobbyHandler) JoinLobby() gin.HandlerFunc {
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
		joinCode := c.PostForm("join_code")
		gameId, err := h.lobbyRepo.GameIdLobby(c.Request.Context(), joinCode)
		if err != nil {
			// the join code is removed when the game starts: an unknown code may also be a game already started
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "This join code doesn't exist or the game has already started"})
			return
		}
		kicked, err := h.lobbyRepo.WasKicked(c.Request.Context(), userId.(int64), gameId)
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Something went wrong"})
			return
		}
		if kicked {
			c.JSON(http.StatusForbidden,
				gin.H{"message": "You have been kicked out of this game"})
			return
		}
		lobby, err := h.lobbyRepo.ActualLobbyGameId(c.Request.Context(),
			userId.(int64), gameId)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "Can't find the actual lobby"})
			return
		}
		if lobby != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "You can't join the lobby twice"})
			return
		}
		playerNb, err := h.lobbyRepo.LobbyPlayerNb(c.Request.Context(), gameId)
		if err != nil{
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return
		}
		gamePlayerNb, err := h.gameRepo.PlayerNb(c.Request.Context(), gameId)
		if playerNb + 1 > int(gamePlayerNb) {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The lobby is full"})
			return
		}
		lobbyObj := &games.Lobby{UserId: userId.(int64),
								GameId: gameId,
								JoinCode: nil,
								IsCreator: false}
		err = h.lobbyRepo.JoinLobby(c.Request.Context(), lobbyObj)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success"})
	}
}

// CurrentLobby returns the game_id of the lobby the user is currently in
// game_id is null if the user isn't in any lobby
func (h *LobbyHandler) CurrentLobby() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		lobby, err := h.lobbyRepo.ActualLobby(c.Request.Context(),
			userId.(int64))
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Something went wrong"})
			return
		}
		if lobby == nil {
			c.JSON(http.StatusOK,
				gin.H{"message": "success", "game_id": nil})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success", "game_id": lobby.GameId})
	}
}

func (h *LobbyHandler) LeaveLobby() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "failed"})
			return
		}
		lobby, err := h.lobbyRepo.ActualLobby(c.Request.Context(),
			userId.(int64));
		if err != nil {
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "failed"})
			return
		}
		if lobby == nil {
			// not in a lobby: without this check lobby.ID below panics
			c.JSON(http.StatusNotFound,
				gin.H{"message": "You are not in a lobby"})
			return
		}
		err = h.lobbyRepo.SetLeftAt(c.Request.Context(), lobby.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return
		}
		lobbies, err := h.lobbyRepo.AllPlayers(c.Request.Context(), lobby.GameId)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return
		}
		err = h.lobbyRepo.SetJoinCode(c.Request.Context(),
			lobby.ID, nil)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return
		}
		err = h.lobbyRepo.SetIsCreator(c.Request.Context(),
			lobby.ID, false)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return
		}
		for _, newLobby := range lobbies {
			// only the host hands the host role and the join code over (otherwise any player leaving made someone else host too)
			if lobby.IsCreator && newLobby.UserId != userId.(int64) {
				err = h.lobbyRepo.SetIsCreator(c.Request.Context(),
					newLobby.ID, true)
				if err != nil {
					c.JSON(http.StatusUnauthorized,
						gin.H{"message": "failed"})
					return
				}
				err = h.lobbyRepo.SetJoinCode(c.Request.Context(),
					newLobby.ID, lobby.JoinCode)
				if err != nil {
					c.JSON(http.StatusUnauthorized,
						gin.H{"message": "failed"})
					return
				}
				break
			}
		}
		c.JSON(http.StatusOK, 
			gin.H{"message": "success"})
	}
}

func (h *LobbyHandler) RemovePlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		username := c.Param("username")
		user, err := h.userRepo.UserByName(c.Request.Context(), username)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": "can't get the user by name"})
			return
		}
		if user.ID == userId.(int64) {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": "you can't kick yourself out of the lobby"})
			return
		}
		lobby, err := h.lobbyRepo.ActualLobby(c.Request.Context(), user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": "failed"})
			return
		}
		if lobby == nil {
			c.JSON(http.StatusNotFound,
				gin.H{"message": "this player is not in a lobby"})
			return
		}
		isCreator, err := h.lobbyRepo.IsCreator(c.Request.Context(),
			userId.(int64), lobby.GameId)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": "can't get the lobby"})
			return
		}
		if !isCreator {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "you're not allowed to remove the player"})
			return
		}
		// Kick (and not SetLeftAt): the kicked player won't be able to join this game again
		err = h.lobbyRepo.Kick(c.Request.Context(), lobby.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return
		}
		c.JSON(http.StatusOK, 
			gin.H{"message": "success"})
	}
}

func (h *LobbyHandler) AllPlayers() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		lobby, err := h.lobbyRepo.ActualLobby(c.Request.Context(),
			userId.(int64))
		if err != nil || lobby == nil {
			c.JSON(http.StatusNotFound,
				gin.H{"message": "You are not in a lobby"})
			return 
		}
		lobbies, err := h.lobbyRepo.AllPlayers(c.Request.Context(),
			lobby.GameId)
		if err != nil {
			c.JSON(http.StatusNotFound,
				gin.H{"message": "There's none in the lobby"})
			return 
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success", "lobbies": lobbies})
	}
}

func (h *LobbyHandler) NbGames() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("userid"), 10, 64)
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "An error occured. Please try again."})
			return
		} 
		nbGames, err := h.lobbyRepo.NbGames(c.Request.Context(), userId)
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "An error occured. Please try again."})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success", "nb_games": nbGames})
	}
}

func (h *LobbyHandler) Historic() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("userid"), 10, 64)
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusBadRequest, 
				gin.H{"message": "An error occured. Please try again."})
			return
		} 
		historicGames, err := h.lobbyRepo.Historic(c.Request.Context(), userId)
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusInternalServerError, 
				gin.H{"message": "An error occured. Please try again."})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success", "historic_games": historicGames})
	}
}