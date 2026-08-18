package handlers

import (
	// "strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"backend/internal/models/games"
	"backend/internal/repositories"
)

type LobbyHandler struct {
	gameRepo *repositories.GameRepository
	lobbyRepo *repositories.LobbyRepository
	userRepo *repositories.UserRepository
}

func NewLobbyHandler(gameRepo *repositories.GameRepository,
	lobbyRepo *repositories.LobbyRepository,
	userRepo *repositories.UserRepository) *LobbyHandler {
	return &LobbyHandler{gameRepo: gameRepo, lobbyRepo: lobbyRepo,
		userRepo: userRepo}
}

// id, err := strconv.ParseInt(c.Param("Join"), 10, 64)
func (h *LobbyHandler) JoinLobby() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		joinCode := c.PostForm("join_code")
		gameId, err := h.lobbyRepo.GameIdLobby(c.Request.Context(), joinCode)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return
		}
		lobby, err := h.lobbyRepo.LobbyUserGameId(c.Request.Context(),
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
								GameId: gameId}
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

func (h *LobbyHandler) LeaveLobby() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "failed"})
			return
		}
		// gameId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, 
		// 		gin.H{"message": "invalid id"})
		// 	return;
		// }
		lobby, err := h.lobbyRepo.LobbyUser(c.Request.Context(),
			userId.(int64));
		if err != nil {
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "failed"})
			return
		}
		err = h.lobbyRepo.RemoveLobby(c.Request.Context(), lobby)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": err})
			return
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
				gin.H{"message": "cant't get the user by name"})
			return
		}
		lobby, err := h.lobbyRepo.ActualLobby(c.Request.Context(), user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": "failed"})
			return
		}
		joinCode, err := h.lobbyRepo.IsAllowedToRemove(c.Request.Context(),
			userId.(int64), lobby.GameId)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": "cant't get the lobby"})
			return
		}
		if len(joinCode) == 0 {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "you're not allowed to remove the player"})
			return
		}
		err = h.lobbyRepo.RemoveLobby(c.Request.Context(), lobby)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "failed"})
			return
		}
		c.JSON(http.StatusOK, 
			gin.H{"message": "success"})
	}
}