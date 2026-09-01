package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"encoding/json"
	"crypto/subtle"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"backend/internal/models/users"
	"backend/internal/repositories"
	"backend/internal/tools/code_generation"
	"backend/internal/tools/cookie"
	"backend/internal/tools/encryption"
	"backend/internal/tools/mailing"
	"backend/internal/tools/token_generation"
)

/**
 * README:
 * Handlers make the bridge between the router and the repository functions.
 * It's important to split different features into different files to improve visibility
 * and maintainability.
 *
 * Basically, the handlers call the repository functions, manage any errors,
 * and return the result to the frontend (in JSON, plain text, etc.).
 */
 type UserHandler struct {
	userRepo     *repositories.UserRepository
	aTokenRepo   *repositories.ActivationTokenRepository
	tokenHandler *TokenHandler
	googleConf   *oauth2.Config
}

func NewUserHandler(
	userRepo *repositories.UserRepository,
	aTokenRepo *repositories.ActivationTokenRepository,
	tokenHandler *TokenHandler,
	googleConf *oauth2.Config,
) *UserHandler {
	return &UserHandler{
		userRepo:     userRepo,
		aTokenRepo:   aTokenRepo,
		tokenHandler: tokenHandler,
		googleConf:   googleConf,
	}
}

func (h *UserHandler) createActivationToken(c *gin.Context, expiresAt time.Time,
	userId int64) error {
	tokenObj, tokenErr := token_generation.GenerateToken(userId, expiresAt,
		[]byte(os.Getenv("ACTIVATION_JWT_SECRET")))
	if tokenErr != nil {
        c.JSON(http.StatusUnauthorized,
            gin.H{"message": "Send a new activation link"})
		return tokenErr
	}
	err := h.aTokenRepo.DeleteActivationToken(c.Request.Context(), userId)
	if err != nil {
        c.JSON(http.StatusInternalServerError,
            gin.H{"message": "Send a new activation link"})
		return err
	}
	err = h.aTokenRepo.CreateAToken(c.Request.Context(), tokenObj, expiresAt)
	if err != nil {
        fmt.Printf("Token creation error message : %s\n", err)
        c.JSON(http.StatusInternalServerError,
            gin.H{"message": "Send a new activation link"})
		return err
	}
	email := c.PostForm("email")
	if len(email) == 0 {
		emailString, err := h.userRepo.Email(c.Request.Context(), userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "can't find the email"})
			return err
		}
		err = mailing.SendActivationEmail(*emailString, tokenObj.Token);
	} else {
		err = mailing.SendActivationEmail(email, tokenObj.Token);
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Send a new activation link"})
		fmt.Printf("\n%s\n", err)
		return err
	}
	return nil
}

func (h *UserHandler) jwtGenereation() (*tokens.ActivationToken, bool){
	expiresAt := time.Now().Add(15 * time.Minute)
	tokenObj, tokenErr := token_generation.GenerateToken(userId, expiresAt,
		[]byte(os.Getenv("JWT_SECRET")))
	if tokenErr != nil {
        c.JSON(http.StatusInternalServerError,
            gin.H{"message": "token creation went wrong"})
		return nil, false
	}
	encodedKey, rTokenErr := h.tokenHandler.CreateRefreshToken(c.Request.Context(), userId)
	if rTokenErr != nil {
        c.JSON(http.StatusInternalServerError,
            gin.H{"message": "Try to login again"})
		return nil, false
	}
	cookie.Set(c, "/api/auth", "refresh_token",
		encodedKey, 7*24*time.Hour, true)
	return tokenObj, true
}

func (h *UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		hashpassword, hashErr := encryption.HashPassword(c.PostForm("password"))
		if hashErr != nil {
            c.JSON(http.StatusInternalServerError,
                gin.H{"message": "Something went wrong"})
			return
		}
		login := c.PostForm("login")
		username := c.PostForm("username")
		email := c.PostForm("email")
		if len(login) == 0 || len(username) == 0 ||
			len(email) == 0 {
		    c.JSON(http.StatusBadRequest,
                gin.H{"message": "failed"})
			return
		}
		if len(username) == 0 {
		    username = "player"	
		}
		userObj := &users.User{Login: login,
			Password:    &hashpassword,
			Username:    username,
			Email:       email,
			Is_verified: false}

		user, err := h.userRepo.CreateUser(c.Request.Context(), userObj)
		if err != nil {
            c.JSON(http.StatusInternalServerError,
                gin.H{"message": "The account creation failed"})
			return
		} else {
			expiresAt := time.Now().Add(24 * time.Hour)
			if err := h.createActivationToken(c, expiresAt, user.ID); err != nil {
				return
			}
            c.JSON(http.StatusCreated,
                gin.H{"message": "The activation link has been sent on your email"})
		}
	}
}

func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		state, userId, err := h.userRepo.LoginUser(c.Request.Context(),
			c.PostForm("login"), c.PostForm("password"))
		if err != nil && userId != 0 {
			expiresAt := time.Now().Add(24 * time.Hour)
			if err := h.createActivationToken(c, expiresAt, userId); err != nil {
				return
			}
			c.JSON(http.StatusOK,
                gin.H{"message": "The activation link has been sent on your email"})
            return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"state": "failed"})
			return
		}
		if !state {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "Try to login again"})
			return
		} else {
			// generate the jwt token here
			tokenObj, ok := jwtGeneration()
			if !ok {
				return
			}
            c.JSON(http.StatusOK,
                gin.H{"Authorization": tokenObj.Token})
		}
	}
}

func (h *UserHandler) GoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := code_generation.RandGenerate(32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, 
				gin.H{"message": "failed"})
			return
		}
		cookie.Set(c, "/api/auth/google", "oauth_state",
		state, 10*time.Minute, false)
		url := h.googleConf.AuthCodeURL(state)
		c.Redirect(http.StatusFound, url)
	}
}

func (h *UserHandler) GoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		googleError := c.Query("error")
		googleState := c.Query("state")
		userState, err := c.Cookie("oauth_state")
		if err != nil || googleError{
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "failed"})
			return 
		}
		if subtle.ConstantTimeCompare([]byte(userState), []byte(googleState)) != 1 {
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "failed"})
			return 
		}
		code := c.Query("code")
		cookie.Clear(c, "/api/auth/google", "oauth_state") //destroy the cookie
		token, err := h.googleConf.Exchange(c, code) //google give us a token

		// use the token to get the user info
		client := h.googleConf.Client(c, token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "failed"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "failed"})
			return
		}

		var userInfo users.GoogleUserInfo
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
			return
		}
		// perform the last checks
		if !userInfo.EmailVerified {
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "failed"})
			return
		}
		user, err := h.userRepo.UserByEmail(c.Request.Context(), userInfo.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, 
				gin.H{"message": "failed"})
			return
		}
		if user == nil {
			// we should create the user
		} else {
			if user.GoogleId == nil {
				err := h.setGoogleId(c, user.ID, ???)
				if err != nil {
					c.JSON(http.StatusInternalServerError,
						gin.H{"message": "token creation went wrong"})
					return
				}
			}
		}
		// generate the jwt token here
		tokenObj, ok := jwtGeneration()
		if !ok {
			return
		}
		c.JSON(http.StatusOK, 
			gin.H{"message": "success", "Authorization": tokenObj.Token})
	}
}

// Getters
func (h *UserHandler) User() gin.HandlerFunc {
	//return almost all user infos
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "The userId is empty"})
            return
        }
		userData, err := h.userRepo.User(c.Request.Context(), userId.(int64))
		if err != nil {
		    c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
            return
        }
		if userData == nil {
            c.JSON(http.StatusNotFound,
                gin.H{"message": "The user doesn't exist"})
            return
		} else {
            c.JSON(http.StatusOK,
                gin.H{"message": "success", "user":  userData})
            return
		}
	}
}
func (h *UserHandler) Name() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "The user id is empty"})
		}
		name, err := h.userRepo.Name(c.Request.Context(), userId.(int64))
		if err != nil {
		    c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		if name == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "The user doesn't have a name"})
			return 
		} else {
            c.JSON(http.StatusOK,
                gin.H{"message": "success", "name":  name})
		}
	}
}
func (h *UserHandler) GetLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "The user id is empty"})
		}
		login, err := h.userRepo.Login(c.Request.Context(), userId.(int64))
		if err != nil {
		    c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		if login == nil {
		    c.JSON(http.StatusNotFound, gin.H{"message": "The user doesn't have a login"})
			return
		} else {
            c.JSON(http.StatusOK, gin.H{"message": "success", "login":  login})
		}
	}
}
func (h *UserHandler) Email() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "The user id is empty"})
			return
		}
		email, err := h.userRepo.Email(c.Request.Context(), userId.(int64))
		if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
			return
		}
		if email == nil {
		    c.JSON(http.StatusNotFound, gin.H{"message": "The user doesn't have a email"})
			return
		} else {
            c.JSON(http.StatusOK, gin.H{"message": "success", "email":  email})
		}
	}
}

func (h *UserHandler) UserByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		user, err := h.userRepo.UserByName(c.Request.Context(), username)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": "failed"})
			return		
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success", "user": user})
	}
}

// Setters
func (h *UserHandler) SetName() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "The user id is empty"})
			return
		}
		state, err := h.userRepo.SetName(c.Request.Context(), userId.(int64), c.PostForm("name"))
		if err != nil {
		    c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
			return
		}
		if !state {
		    c.JSON(http.StatusOK, gin.H{"message": "The username has not been modified"})
			return
		} else {
            c.JSON(http.StatusOK, gin.H{"message": "success"})
		}
	}
}
func (h *UserHandler) setGoogleId(c *gin.Context, userId int64, googleId *string) bool {
	state, err := h.userRepo.SetGoogleId(c.Request.Context(), userId, googleId)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
		return false
	}
	if !state {
	    c.JSON(http.StatusOK, gin.H{"message": "The googleId has not been modified"})
		return false
	} else {
		return true
	}
}
func (h *UserHandler) SetLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "The user id is empty"})
			return
		}
		state, err := h.userRepo.SetLogin(c.Request.Context(), userId.(int64), c.PostForm("login"))
		if err != nil {
		    c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
			return
        }
		if !state {
		    c.JSON(http.StatusOK, gin.H{"message": "The login has not been modified"})
			return
		} else {
            c.JSON(http.StatusOK, gin.H{"message": "success"})
		}
	}
}
func (h *UserHandler) SetAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "The user id is empty"})
			return
		}
		state, err := h.userRepo.SetAvatar(c.Request.Context(), userId.(int64), c.PostForm("avatar"))
		if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
			return
		}
        if !state {
            c.JSON(http.StatusOK, gin.H{"message": "The avatar has not been modified"})
			return
		} else {
             c.JSON(http.StatusOK, gin.H{"message": "success"})
        }
	}
}
