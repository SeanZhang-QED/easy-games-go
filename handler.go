package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getGameHandler(w http.ResponseWriter, r *http.Request) {
	
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

	urlParams := r.URL.Query()

	var games []Game
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
