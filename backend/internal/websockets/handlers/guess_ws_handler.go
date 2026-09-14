package ws_handlers

import (
	"backend/internal/repositories"
	"context"
	"backend/internal/handlers"
	"backend/internal/models/guesses"
	"backend/internal/tools/haversine"
	"errors"
	"github.com/sams96/rgeo"
	"log"
)

type GuessWsHandler struct {
	guessRepo *repositories.GuessRepository
	roundRepo *repositories.RoundRepository
	lobbyRepo *repositories.LobbyRepository;
	gameRepo *repositories.GameRepository;
	reverser *rgeo.Rgeo
}

func NewGuessWsHandler(guessRepo *repositories.GuessRepository, 
	roundRepo *repositories.RoundRepository, lobbyRepo *repositories.LobbyRepository, 
	gameRepo *repositories.GameRepository, reverser *rgeo.Rgeo) *GuessWsHandler {
	return &GuessWsHandler{guessRepo: guessRepo,
		roundRepo: roundRepo, lobbyRepo: lobbyRepo, gameRepo: gameRepo,
		reverser: reverser}
}

type RoundInfos struct {
	Distance float64 `json:"distance"`
	Score int64 `json:"score"`
	CountryFound bool `json:"country_found"`
	BaseCC string `json:"base_cc"`
	GuessedCC string `json:"guessed_cc"`
	Streak int64 `json:"streak"`
	TimeTaken time.time `json:"time_taken"`
}

type Guess struct {
	Username string `json:"username"`
	RoundNb int64 `json:"round_nb"`
	RoundInfos RoundInfos `json:"round_infos"`
	TotalScore int64 `json:"total_score"`
}

//TODO comment renvoyer tous les score de tous les joueurs a chaque joueur ? fautil renvoyer une map par username? ou renvoyer les structure une par une?

//ADDGUESS A ADAPTER POUR LES WEBSOCKETS
func (h *GuessWsHandler) AddGuess(ctx context.Context, serializedData) (*guesses.Guess, bool) {
	// deserialiser
	var guessData Guess//map[string]interface{}
	err := json.Unmarshal(serializedData.Payload, &guessData)
	if err != nil {
		log.Printf("error : %s\n", err)
		return nil, false
	}
	userID := serializedData.Client.userId
	gameID := serializedData.Client.gameId
	if guess.Lat == nil || guess.Long == nil {
    	// c.JSON(http.StatusBadRequest, gin.H{"error": "latitude and longitude must be provided"})
    	return nil, false
	}
	
	if guess.RoundNb == 0 {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "round_nb must be provided"})
		return nil, false
	} // round_nb permet de recuperer le round_id en coordonnant avec le user_id (et checker si la game est en cours, et si le user a deja guess pour ce round, et si le round est bien celui de la game en cours)

	if userID == 0 {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "user_id must be provided"})
		return nil, false // a coordonner avec le round_nb
	}


	// // TODO do I still need this part??
	// // check if the user is in an active lobby (to get the game_id and check if the round_nb is valid for this game)
	// lobby, err := h.lobbyRepo.ActualLobby(ctx, userID)
	// if err != nil {
	// 	log.Printf("Error getting active lobby for user %d: %v", userID, err)
	// 	// c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get active lobby"})
	// 	return nil, false
	// }
	// if lobby == nil {
	// 	// c.JSON(http.StatusBadRequest, gin.H{"error": "user is not in an active lobby"})
	// 	return nil, false
	// }



	// check if the round_nb exists in table rounds for this specific game id (if it is the next round to be played)
	guess.RoundID, err = h.roundRepo.GetIDByGameIDAndRoundNb(ctx, gameID, guess.RoundNb)
	if err != nil {
		log.Printf("Error getting round id for game %d and round_nb %d: %v", gameID, guess.RoundNb, err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get round id"})
		return nil, false
	}

	// check if the game is finished (if true, then the user cannot guess anymore) OR check if round_nb is less than round_total in the game table
	finished, err := h.gameRepo.Finished(ctx, gameID)
	if err != nil {
		log.Printf("Error checking if game %d is finished: %v", gameID, err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "could not check game status"})
		return nil, false
	}
	if finished {
		log.Printf("Game %d is already finished, user %d cannot guess", gameID, userID)
		// c.JSON(http.StatusBadRequest, gin.H{"error": "this game is already finished"})
		return nil, false
	}

	// TODO check if the player has not left the lobby

	// round_nb should be the same as the current round (last round created for the game), else the user is trying to guess for a round that is not the current one (or has already been played)
	nextRoundNb, err := h.roundRepo.GetNextRoundNb(ctx, gameID)
	if err != nil {
		log.Printf("Error getting current round number for game %d: %v", gameID, err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "could not verify round number"})
		return nil, false
	}
	if guess.RoundNb != nextRoundNb - 1 {
		log.Printf("User %d is trying to guess for round %d, but the current round for game %d is %d",
		 userID, guess.RoundNb, gameID, nextRoundNb - 1)
		// c.JSON(http.StatusBadRequest, gin.H{"error": "round_nb does not match the current round for this game"})
		return nil, false
	}

	// get locationID using roundID (to get base coordinates and base country code)
	locationID, err := h.guessRepo.GetLocationIDByRoundID(ctx, guess.RoundID)
	if err != nil {
		log.Printf("Error getting location ID for user %d at round %d: %v", userID, guess.RoundID, err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get location ID for round"})
		return nil, false
	}
	// get base coordinates using locationID
	log.Printf("Location ID for round %d: %d", guess.RoundID, locationID) // TODO: remove this log after testing
	baseLat, baseLong, err := h.roundRepo.GetCoordinatesByLocationID(ctx, locationID)
	if err != nil {
		log.Printf("Error getting base coordinates for round %d: %v", guess.RoundID, err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get base coordinates for round"})
		return nil, false
	}
	log.Printf("Base coordinates for round %d: (%f, %f)", guess.RoundID, baseLat, baseLong) // TODO: remove this log after testing
	log.Printf("Guess coordinates: (%f, %f)", *guess.Lat, *guess.Long) // TODO: remove this log after testing
	//get base cc using locationID
	baseCC, err := h.roundRepo.GetCCByLocationID(ctx, locationID)
	if err != nil {
		log.Printf("Error getting country code for location %d: %v", locationID, err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get country code for location"})
		return nil, false
	}
	log.Printf("Country code for location %d: %s", locationID, baseCC) // TODO: remove this log after testing

	//get cc using guess coordinates and package georeverse
	loc, err := h.reverser.ReverseGeocode([]float64{*guess.Long, *guess.Lat})
	if err != nil {
		log.Printf("Error reverse geocoding guess coordinates (%f, %f): %v", *guess.Lat, *guess.Long, err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "could not reverse geocode guess coordinates"})
		return nil, false
	}
	guess.CC = loc.CountryCode2
	log.Printf("Country code for this guess: %s", guess.CC) // TODO: remove this log after testing
	if guess.CC == baseCC {
		guess.CountryFound = true
	} 

	// Calculate distance and score before persisting
	guess.Distance = handlers.CalculateDistance(haversine.Coordinates{Lat: baseLat, Lon: baseLong}, haversine.Coordinates{Lat: *guess.Lat, Lon: *guess.Long})
	handlers.CalculateScore(&guess)

	createdGuess, err := h.guessRepo.Create(ctx, &guess)
	if err != nil {
		log.Printf("Error creating guess: %v", err)
		// to check if it is a unique constraint violation (means user already guessed for this round)
		if errors.Is(err, repositories.ErrGuessNotFound) || err.Error() == "pq: duplicate key value violates unique constraint" {
			// c.JSON(http.StatusConflict, gin.H{"error": "user has already submitted a guess for this round"})
			return nil, false
		}
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create guess"})
		return nil, false
	}

	// c.JSON(http.StatusCreated, createdGuess)
	for client := range h.clients[gameId] {
		serializedResponse, err := json.Marshal(usersLobby) // TODO replace usersLobby with the actual response data you want to send to clients
		if err != nil {
			log.Printf("error : %s\n", err)
			return 
		}
		client.send <- serializedResponse
	return createdGuess, true
}