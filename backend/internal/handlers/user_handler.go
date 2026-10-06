package handlers

import (
	"strings"
	"crypto/sha256"
	"log"
	"fmt"
	"net/http"
	"os"
	"time"
	"errors"
	"encoding/json"
	"crypto/subtle"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"path/filepath"
	"github.com/go-playground/validator/v10"

	"backend/internal/models/users"
	"backend/internal/repositories"
	"backend/internal/tools/code_generation"
	"backend/internal/tools/cookie"
	"backend/internal/tools/encryption"
	"backend/internal/tools/mailing"
	"backend/internal/tools/token_generation"
	"backend/internal/models/users/tokens"
	"backend/internal/tools/rank_calculation"
	"backend/internal/tools/validation"
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
	rTokenRepo   *repositories.RefreshTokenRepository
	tokenHandler *TokenHandler
	googleConf   *oauth2.Config
	validate	 *validator.Validate
}

func NewUserHandler(
	userRepo *repositories.UserRepository,
	aTokenRepo *repositories.ActivationTokenRepository,
	rTokenRepo *repositories.RefreshTokenRepository,
	tokenHandler *TokenHandler,
	googleConf *oauth2.Config,
) *UserHandler {

	validate := validator.New()
	validation.RegisterCustomRules(validate)

	return &UserHandler{
		userRepo:     userRepo,
		aTokenRepo:   aTokenRepo,
		rTokenRepo:   rTokenRepo,
		tokenHandler: tokenHandler,
		googleConf:   googleConf,
		validate:	  validate,
	}
}

func (h *UserHandler) createActivationToken(c *gin.Context, expiresAt time.Time,
	userId int64) error {
	tokenObj, tokenErr := token_generation.GenerateToken(userId, expiresAt,
		[]byte(os.Getenv("ACTIVATION_JWT_SECRET")))
	if tokenErr != nil {
        c.JSON(http.StatusUnauthorized,
			gin.H{"message": "The activation link sending went wrong. Please try again."})
		return tokenErr
	}
	err := h.aTokenRepo.DeleteActivationToken(c.Request.Context(), userId)
	if err != nil {
        c.JSON(http.StatusInternalServerError,
			gin.H{"message": "The activation link sending went wrong. Please try again."})
		return err
	}
	err = h.aTokenRepo.CreateAToken(c.Request.Context(), tokenObj, expiresAt)
	if err != nil {
        fmt.Printf("Token creation error message : %s\n", err)
        c.JSON(http.StatusInternalServerError,
			gin.H{"message": "The activation link sending went wrong. Please try again."})
		return err
	}
	email := c.PostForm("email")
	if len(email) == 0 {
		emailString, err := h.userRepo.Email(c.Request.Context(), userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Can't find your email. Please try again."})
			return err
		}
		err = mailing.SendActivationEmail(*emailString, tokenObj.Token);
	} else {
		err = mailing.SendActivationEmail(email, tokenObj.Token);
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "The activation link sending went wrong. Please try again."})
		fmt.Printf("\n%s\n", err)
		return err
	}
	return nil
}

func (h *UserHandler) jwtGeneration(c *gin.Context, userId int64) (*tokens.ActivationToken, bool) {
	expiresAt := time.Now().Add(15 * time.Minute)
	tokenObj, tokenErr := token_generation.GenerateToken(userId, expiresAt,
		[]byte(os.Getenv("JWT_SECRET")))
	if tokenErr != nil {
        c.JSON(http.StatusInternalServerError,
            gin.H{"message": "token creation went wrong."})
		return nil, false
	}
	encodedKey, rTokenErr := h.tokenHandler.CreateRefreshToken(c.Request.Context(), userId)
	if rTokenErr != nil {
        c.JSON(http.StatusInternalServerError,
            gin.H{"message": "Try to login again."})
		return nil, false
	}
	cookie.Set(c, "/api/auth", "refresh_token",
		encodedKey, 7*24*time.Hour, true)
	cookie.Set(c, "/api", "Authorization",
		tokenObj.Token, 15 * time.Minute, true)
	return tokenObj, true
}

type CreateUserValidation struct {
	Email   	string `validate:"required,max=254,email"`
	Login	  	string `validate:"required,min=3,max=30,acceptChar"`
	Username	string `validate:"required,min=3,max=30,acceptChar"`
	Password	string `validate:"required,strongpwd"`
}

func (h *UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		input := CreateUserValidation{
			Email:		strings.ToLower(strings.TrimSpace(c.PostForm("email"))),
			Login:		strings.ToLower(strings.TrimSpace(c.PostForm("login"))),
			Username:	strings.ToLower(strings.TrimSpace(c.PostForm("username"))),
			Password:	c.PostForm("password"),
		}
		err := h.validate.Struct(input)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": validation.FormatErrors(err)})
			return
		}
		hashpassword, hashErr := encryption.HashPassword(input.Password)
		if hashErr != nil {
            c.JSON(http.StatusInternalServerError,
                gin.H{"message": "Something went wrong. Please try again."})
			return
		}
		existingLogin, err := h.userRepo.ExistingLogin(c.Request.Context(), input.Login)
		if err != nil {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "Something went wrong. Please try again."})
			return
		}
		if existingLogin != nil && *existingLogin != "" {
            c.JSON(http.StatusBadRequest,
                gin.H{"message": "This login or this username are unavailable. Please try another one."})
			return
		}
		existingUsername, err := h.userRepo.ExistingName(c.Request.Context(), input.Username)
		if err != nil {
            c.JSON(http.StatusUnauthorized,
                gin.H{"message": "Something went wrong. Please try again."})
			return
		}
		if existingUsername {
            c.JSON(http.StatusBadRequest,
                gin.H{"message": "This username is unavailable. Please try another one."})
			return
		}		 
		userObj := &users.User{Login: &input.Login,
			Password:    &hashpassword,
			Username:    input.Username,
			Email:       input.Email,
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

type LoginValidation struct {
	Login	  	string `validate:"required"`
	Password	string `validate:"required"`
}

func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		input := LoginValidation {
			Login:		strings.ToLower(strings.TrimSpace(c.PostForm("login"))),
			Password:	c.PostForm("password"),
		}
		err := h.validate.Struct(input)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": validation.FormatErrors(err)})
			return
		}
		state, userId, err := h.userRepo.LoginUser(c.Request.Context(),
			input.Login, input.Password)
		if err != nil && userId != 0 {
			expiresAt := time.Now().Add(24 * time.Hour)
			if err := h.createActivationToken(c, expiresAt, userId); err != nil {
				return
			}
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusUnauthorized,
                gin.H{"message": "The activation link has been resent on your email"})
            return
		} else if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "Invalid login ID or password. Please try again."})
			return
		}
		if !state {
            c.JSON(http.StatusBadRequest,
                gin.H{"message": "Invalid login ID or password. Please try again."})
			return
		} else {
			// generate the jwt token in the cookies
			_, ok := h.jwtGeneration(c, userId)
			if !ok {
				return
			}
            c.JSON(http.StatusOK,
                gin.H{"message": "success"})
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

func (h *UserHandler) fetchGoogleUser(c *gin.Context, code string) (*users.GoogleUserInfo, error) {
	token, err := h.googleConf.Exchange(c, code)
	if err != nil {
		return nil, err
	}
	resp, err := h.googleConf.Client(c, token).Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo: status %d", resp.StatusCode)
	}
	var info users.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (h *UserHandler) resolveGoogleUser(c *gin.Context, info *users.GoogleUserInfo) (*users.User, error) {
	user, err := h.userRepo.UserByGoogleId(c.Request.Context(), &info.Sub)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	user, err = h.userRepo.UserByEmail(c.Request.Context(), info.Email)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, err
	}
	if user != nil {
		ok := h.setGoogleId(c, user.ID, &info.Sub)
		if !ok {
			log.Printf("error : failed\n")
			return nil, errors.New("failed")
		}
		return user, nil
	}
	// change the username if it already exist
	checkUsername := info.Name
	for i := 0; i < 50; i++ {
		existingName, err := h.userRepo.ExistingName(c.Request.Context(), checkUsername)
		if err != nil  {
			log.Printf("error : %s\n", err)
			return nil, errors.New("failed")
		}
		if existingName == false {
			break 
		}
		code, err := code_generation.RandGenerate(4)
		if err != nil {
			log.Printf("error : %s\n", err)
			return nil, errors.New("failed")
		}
		checkUsername = fmt.Sprintf("%s_%s", info.Name, code)
	}
	userObj := &users.User{Login: nil,
			Password:    nil,
			Username:    strings.ToLower(checkUsername),
			Email:       strings.ToLower(info.Email),
			Is_verified: true,
			GoogleId: 	 &info.Sub,}
	user, err = h.userRepo.CreateUser(c.Request.Context(), userObj)
	if err != nil {
		log.Printf("error : The account creation failed\n")
		return nil, errors.New("failed")
	}
	return user, nil
}

func (h *UserHandler) GoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		googleState := c.Query("state")
		userState, err := c.Cookie("oauth_state")
		queryError := c.Query("error")
		loginUrl := fmt.Sprintf("%s/login", os.Getenv("PUBLIC_URL"))
		homeUrl := fmt.Sprintf("%s/home", os.Getenv("PUBLIC_URL"))
		if err != nil || queryError != "" {
			if err != nil {
				log.Printf("error : %s\n", err)
			} else if queryError != "" {
				log.Printf("error : %s\n", queryError)
			}
			c.Redirect(http.StatusFound, loginUrl)
			return
		}
		if subtle.ConstantTimeCompare([]byte(userState), []byte(googleState)) != 1 {
			c.Redirect(http.StatusFound, loginUrl)
			return
		}
		cookie.Clear(c, "/api/auth/google", "oauth_state")

		userInfo, err := h.fetchGoogleUser(c, c.Query("code"))
		if err != nil {
			c.Redirect(http.StatusFound, loginUrl)
			return 
		}
		if !userInfo.EmailVerified {
			c.Redirect(http.StatusFound, loginUrl)
			return
		}
		user, err := h.resolveGoogleUser(c, userInfo)
		if err != nil {
			c.Redirect(http.StatusFound, loginUrl)
			return
		}
		_, ok := h.jwtGeneration(c, user.ID)
		if !ok {
			c.Redirect(http.StatusFound, loginUrl)
			return
		}
		c.Redirect(http.StatusFound, homeUrl)
	}
}

func (h *UserHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		encodedKey, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You have to login"})
			return
		}
		hashedKey := sha256.Sum256([]byte(encodedKey))
		hashToken := fmt.Sprintf("%x", hashedKey)
		if err := h.rTokenRepo.SetRefreshUsedAt(c.Request.Context(), hashToken); err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "You have to login"})
			return
		}
		cookie.Clear(c, "/api/auth", "refresh_token")
		cookie.Clear(c, "/api", "Authorization")
		c.JSON(http.StatusOK,
			gin.H{"message": "you are now logout"})
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

type SetNameValidation struct {
	Username	string `validate:"required,min=3,max=30,acceptChar"`
}

func (h *UserHandler) SetName() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The user id is empty"})
			return
		}
		input := SetNameValidation {
			Username: strings.ToLower(strings.TrimSpace(c.PostForm("name"))),
		}
		err := h.validate.Struct(input)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"message": validation.FormatErrors(err)})
			return
		}
		state, err := h.userRepo.SetName(c.Request.Context(), userId.(int64), input.Username)
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

func (h *UserHandler) UploadAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5<<20)
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "The user id is empty"})
			return
		}

		oldUser, err := h.userRepo.User(c.Request.Context(), userId.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch user"})
			return
		}

		file, err := c.FormFile("avatar")
		if err != nil {
			if err.Error() == "http: request body too large" {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{"message": "File too large (max 5MB)"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"message": "No file provided"})
			return
		}

		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file type"})
			return
		}

		filename := uuid.New().String() + ext
		savePath := filepath.Join("uploads/avatars", filename)
		avatarURL := "/api/uploads/avatars/" + filename

		err = c.SaveUploadedFile(file, savePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to save file"})
			return
		}

		_, err = h.userRepo.SetAvatar(c.Request.Context(), userId.(int64), avatarURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update avatar"})
			return
		}

		if oldUser.Avatar != nil && strings.HasPrefix(*oldUser.Avatar, "/api/uploads/avatars/") {
			oldFilename := filepath.Base(*oldUser.Avatar)
			oldPath := filepath.Join("uploads/avatars", oldFilename)
			os.Remove(oldPath)
		}

		c.JSON(http.StatusOK, gin.H{"message": "success", "avatar": avatarURL})
	}
}

// XP returns the xp of the connected user (sum of the base_score of all his finished games)
func (h *UserHandler) XPByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		username := c.Param("username")
		user, err := h.userRepo.UserByName(c.Request.Context(), username)
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success", "user_xp": user.XP})
	}
}

// Leaderboard returns the 100 first players ranked by xp
func (h *UserHandler) Leaderboard() gin.HandlerFunc {
	return func(c *gin.Context) {
		leaderboard, err := h.userRepo.Leaderboard(c.Request.Context())
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success", "leaderboard": leaderboard})
	}
}

// UserLeaderboard returns the position, name and xp of the connected user
func (h *UserHandler) UserLeaderboard() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The userId is empty"})
			return
		}
		leaderboard, err := h.userRepo.UserLeaderboard(c.Request.Context(), userId.(int64))
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK,
			gin.H{"message": "success", "leaderboard": leaderboard})
	}
}

func (h *UserHandler) SetRank(xp int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized,
				gin.H{"message": "The user id is empty"})
			return
		}
		rank := rank_calculation.CalculateRank(xp)

		_, err := h.userRepo.SetRank(c.Request.Context(), userId.(int64), rank)
		if err != nil {
			log.Printf("error : %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}