package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SeanZhang-QED/easy-games-go/config"
	"github.com/SeanZhang-QED/easy-games-go/models"
	"github.com/SeanZhang-QED/easy-games-go/twitch"
)

func SearchItemByGameHandler(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin, * means support all the domain
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//Support whitch HTTP headers can be used during the actual request
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//Tells the front-end the type of response will be json
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		// If current requst is a preflight request,
		// only return the header, which has set above.
		return
	}

	// urlParams maps a string key to a list of values.
	//  map[string][]string
	urlParams := r.URL.Query()

	var items map[string][]models.Item
	var err error

	if len(urlParams) == 0 {
		http.Error(w, "Failed to get gameId from this request.", http.StatusBadRequest)
		fmt.Println("Failed to get gameId from this request.")
		return
	} else {
		gameId := urlParams.Get("game_id")
		fmt.Printf("Received a searchItemByGameId request, search for game with id: %v\n", gameId)
		items, err = searchItemsByGameId(gameId)
	}

	// return the error msg
	if err != nil {
		// return msg with http status code
		http.Error(w, "Failed to get result from Twitch API.", http.StatusInternalServerError)
		fmt.Printf("Failed to get result from Twitch API %v. \n", err)
	}

	// Marshal: type To list of Items -> JSON
	itemsJSON, err := json.Marshal(items)

	if err != nil {
		// return msg with http status code
		http.Error(w, "Failed to parse item data from Twitch API.", http.StatusInternalServerError)
		fmt.Printf("Failed to parse item data from Twitch API %v. \n", err)
	}

	w.Write(itemsJSON)
}

func searchItemsByGameId(gameId string) (map[string][]models.Item, error){
	typeToItemSlice := make(map[string][]models.Item)

	var err error
	for _ , itemType := range models.ItemType {
		typeToItemSlice[itemType], err = searchByType(gameId, itemType, config.DEFAULT_SEARCH_LIMIT)
		if err != nil {
			fmt.Printf("Failed to search %s of a specific game from Twitch API %v. \n", itemType, err)
			return nil, err
		}
	}
	return typeToItemSlice, nil
}

func searchByType(gameId string, itemType string, limit int) ([]models.Item, error) {
	var items []models.Item
	var data string
	var err error
	switch itemType {
	case "STREAM":
		data, err = twitch.SearchByGameId(config.STREAM_SEARCH_URL_TEMPLATE, gameId, limit)
	case "CLIP":
		data, err = twitch.SearchByGameId(config.CLIP_SEARCH_URL_TEMPLATE, gameId, limit)
	case "VIDEO":
		data, err = twitch.SearchByGameId(config.VIDEO_SEARCH_URL_TEMPLATE, gameId, limit)
	default:
		fmt.Printf("Unexpected item tpye: %s", itemType)
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	// step 2: convert data into a list of Game struct
	items, err = getItemList(data)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(items); i++ {
		items[i].ItemType = itemType
	}
	return items, nil
}

func getItemList(data string) ([]models.Item, error) {

	var resp models.TwitchItemResponse

	err := json.Unmarshal([]byte(data), &resp)

	if err != nil {
		fmt.Println("Failed to parse data from Twtich response.")
		return nil, err
	}

	return resp.Data, nil
}