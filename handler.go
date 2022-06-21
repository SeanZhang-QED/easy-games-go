package main

import (
	"fmt"
	"net/http"
)

func getGameHandler(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()

	if len(urlParams) == 0 {
		fmt.Println("Received a topGames request.")
	} else {
		gameName := urlParams.Get("game_name")
		fmt.Printf("Received a searchGame request, search for %v\n", gameName)
	}
}
