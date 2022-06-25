package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/SeanZhang-QED/easy-games-go/models"
	"github.com/SeanZhang-QED/easy-games-go/twitch"
	"github.com/SeanZhang-QED/easy-games-go/config"
)

func GetGameHandler(w http.ResponseWriter, r *http.Request) {
	
	//Allow CORS here By * or specific origin, * means support all the domain
    w.Header().Set("Access-Control-Allow-Origin", "*")
	//Support whitch HTTP headers can be used during the actual request
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//Tells the front-end the type of response will be json
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		// If current requst is a preflight request,
        // only return the header, which has set above.
        return
    }

	// urlParams maps a string key to a list of values.
	//  map[string][]string
	urlParams := r.URL.Query()

	var games []models.Game
	var err error

	if len(urlParams) == 0 {
		fmt.Println("Received a topGames request.")
		games, err = topGames(0) 
	} else {
		gameName := urlParams.Get("game_name")
		fmt.Printf("Received a searchGame request, search for %v\n", gameName)
		games, err = searchGame(gameName) 
	}

	// return the error msg
	if err != nil {
		// return msg with http status code
		http.Error(w, "Failed to get result from Twitch API.",http.StatusInternalServerError)
		fmt.Printf("Failed to get result from Twitch API %v. \n", err)
	}

	// Marshal: Game -> JSON
	gamesJSON, err := json.Marshal(games)

	if err != nil {
		// return msg with http status code
		http.Error(w, "Failed to parse game data from Twitch API.",http.StatusInternalServerError)
		fmt.Printf("Failed to parse game data from Twitch API %v. \n", err)
	}

	w.Write(gamesJSON)
}

func topGames(limit int) ([]models.Game, error) {
	// step 1: call searchByName(), send http request to twitch backend
	curLimit := config.DEFAULT_GAME_LIMIT
	if limit != 0 {
		curLimit = limit
	} 
	data, err := twitch.SearchByName(fmt.Sprintf(config.TOP_GAME_URL, curLimit), "")
	if err != nil {
		return nil, err
	}
	
	// step 2: convert data into a list of Game struct
	games, err := getGameList(data)
	if err != nil {
		return nil, err
	}
	return games, nil
}

func searchGame(gameName string) ([]models.Game, error) {
	// step 1: call searchByName(), send http request to twitch backend
	data, err := twitch.SearchByName(config.GAME_SEARCH_URL_TEMPLATE, gameName)
	if err != nil {
		return nil, err
	}

	// step 2: call getGameList(), convert data into a list of Game struct
	games, err := getGameList(data)
	if err != nil {
		return nil, err
	}
	return games, nil
}

func getGameList(data string) ([]models.Game, error) {
	
	resp := models.TwitchGameResponse{}
	
	err := json.Unmarshal([]byte(data), &resp)

	if err != nil {
		fmt.Println("Failed to parse data from Twtich response.")
		return nil, err
	}
	
	return resp.Data, nil
}