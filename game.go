package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Game struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	BoxArtUrl string `json:"box_art_url"`
}

type Pagination struct {
	Cursor string `json:"cursor"`
}

type TwitchGameResponse struct {
	Data []Game `json:"data"`
	Pagination Pagination `json:"pagination"`
}

const (
	TOKEN                    string = "Bearer imh9k8fap8th5jb4w0846bpmexeovg"
	CLIENT_ID                string = "agfj0iyreh884vaf6jp0p0jndj70nl"
	TOP_GAME_URL             string = "https://api.twitch.tv/helix/games/top?first=%v"
	GAME_SEARCH_URL_TEMPLATE string = "https://api.twitch.tv/helix/games"
	DEFAULT_GAME_LIMIT       int    = 20
)

func topGames(limit int) ([]Game, error) {
	// step 1: call searchTwitch(), send http request to twitch backend
	curLimit := DEFAULT_GAME_LIMIT
	if limit != 0 {
		curLimit = limit
	} 
	data, err := searchTwitch(fmt.Sprintf(TOP_GAME_URL, curLimit), "")
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

func searchGame(gameName string) ([]Game, error) {
	// step 1: call searchTwitch(), send http request to twitch backend
	data, err := searchTwitch(GAME_SEARCH_URL_TEMPLATE, gameName)
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

func searchTwitch(url string, gameName string) (string, error) {
	client := http.Client{Timeout: time.Duration(1) * time.Second}

	req, _ := http.NewRequest("GET", url, nil)
	
	if gameName != "" {
		q := req.URL.Query() // Get a copy of the query values.
		q.Add("name", gameName) // Add a new value to the set.
		req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.
	}

	req.Header.Add("Authorization", TOKEN)
	req.Header.Add("Client-Id", CLIENT_ID)

	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to get result from Twitch API.")
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Failed to read result from Twitch response body.")
		return "", err
	}

	bodyString := string(bodyBytes)

	return bodyString, nil
}