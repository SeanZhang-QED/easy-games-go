package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/SeanZhang-QED/easy-games-go/models"
	"github.com/SeanZhang-QED/easy-games-go/session"
	"gopkg.in/mgo.v2"
)

const DEFAULT_GAME_LIMIT int = 3
const DEFAULT_PER_GAME_RECOMMENDATION_LIMIT int = 10
const DEFAULT_TOTAL_RECOMMENDATION_LIMIT int = 20

type RecommendHandler struct {
	session *mgo.Session
}

func NewRecommendHandler(s *mgo.Session) *RecommendHandler {
	return &RecommendHandler{s}
}

func (rh RecommendHandler) Recommendation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one recommendation request")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	// check for login
	_, err := r.Cookie("sessionId")
	if err != nil {
		fmt.Println("Haven't logged in, recommend by default")
		itemMap, err := rh.recommendItemsByDefault()
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			fmt.Println("Error occurs in recommendItemsByDefault() called by Recommendation()")
			return
		}
		mapJSON, err := json.Marshal(itemMap)
		if err != nil {
			// return msg with http status code
			http.Error(w, "Failed to convert items map to JSON", http.StatusInternalServerError)
			fmt.Printf("Failed to convert items map to JSON: %v. \n", err)
			return
		}
		w.Write(mapJSON)
	} else {
		loggedSession, err := session.AlreadyLoggedIn(w, r, rh.session)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		fmt.Printf("User %v logged in, recommend by user's history", loggedSession.Email)
		itemMap, err := rh.recommendItemsByUserHistory(loggedSession.Email)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			fmt.Println("Error occurs in recommendItemsByUserHistory() called by Recommendation()")
			return
		}
		mapJSON, err := json.Marshal(itemMap)
		if err != nil {
			// return msg with http status code
			http.Error(w, "Failed to convert items map to JSON", http.StatusInternalServerError)
			fmt.Printf("Failed to convert items map to JSON: %v. \n", err)
			return
		}
		w.Write(mapJSON)
	}

	fmt.Printf("Get recommend items successfully.\n")
}

// Return a list of Item objects for the given item type(one of [Stream, Video, Clip]).
// Add items are related to the top games provided in the argument
func recommendByTopGames(itemType string, topGames []models.Game) ([]models.Item, error) {
	var recommendedItems []models.Item
	for _, game := range topGames {
		items, err := searchByType(game.Id, itemType, DEFAULT_PER_GAME_RECOMMENDATION_LIMIT)
		if err != nil {
			fmt.Println("Faild to fetch topgame items from twitch.")
			return nil, err
		}
		if len(recommendedItems) == DEFAULT_TOTAL_RECOMMENDATION_LIMIT {
			return recommendedItems, nil
		}
		recommendedItems = append(recommendedItems, items...)
	}
	return recommendedItems, nil
}

// Return a list of Item objects for the given item type(one of [Stream, Video, Clip]).
// All items are related to the items previously favorited by the user.
func recommendByFavoriteHistory(itemType string, favoritedItemIds map[string]models.Void, favoritedGameIds []string) ([]models.Item, error) {
	// step 1: Count the favorite game IDs
	favoriteGameIdByCount := make(map[string]int)
	for _, gameId := range favoritedGameIds {
		if count, ok := favoriteGameIdByCount[gameId]; ok {
			favoriteGameIdByCount[gameId] = count + 1
		} else {
			favoriteGameIdByCount[gameId] = 1
		}
	}

	// step 2: sort the slice of gameId(key) by corresponding count in descending order
	keys := make([]string, 0, len(favoriteGameIdByCount))
	for key := range favoriteGameIdByCount {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return favoriteGameIdByCount[keys[i]] > favoriteGameIdByCount[keys[j]]
	})

	var sortedGameIdByCount []string
	for _, k := range keys {
		fmt.Println(k, ": ", favoriteGameIdByCount[k])
		sortedGameIdByCount = append(sortedGameIdByCount, k)
	}
	if len(sortedGameIdByCount) > DEFAULT_GAME_LIMIT {
		sortedGameIdByCount = sortedGameIdByCount[0:DEFAULT_GAME_LIMIT]
	}

	// step 3: Search Twitch based on the favorite game IDs returned in the last step.
	var recommendedItems []models.Item
	for _, gameId := range sortedGameIdByCount {
		items, err := searchByType(gameId, itemType, DEFAULT_PER_GAME_RECOMMENDATION_LIMIT)
		if err != nil {
			fmt.Println("Faild to fetch topgame items from twitch.")
			return nil, err
		}

		for _, item := range items {
			if len(recommendedItems) == DEFAULT_TOTAL_RECOMMENDATION_LIMIT {
				return recommendedItems, nil
			}
			if _, ok := favoritedItemIds[item.Id]; !ok {
				recommendedItems = append(recommendedItems, item)
			} else {
				continue
			}
		}
	}
	return recommendedItems, nil
}

// Return a map of Item objects for the given item type as key(one of [Stream, Video, Clip]).
// Each key is corresponding to a list of Items objects, each item object is a recommended item
// based on the previous favorite records by the user.
func (rh RecommendHandler) recommendItemsByUserHistory(email string) (map[string][]models.Item, error) {
	recommendedItemMap := make(map[string][]models.Item)
	for _, itemType := range models.ItemType {
		recommendedItemMap[itemType] = []models.Item{}
	}

	favoriteItemIds, err := getFavoriteItemIds(email, rh.session)
	if err != nil {
		fmt.Println("Fail to get users favorite Item history")
		return nil, err
	}
	favoriteGameIds, err := getFavoriteGameIds(email, rh.session)
	if err != nil {
		fmt.Println("Fail to get users favorite Game history")
		return nil, err
	}

	for itemType, gameIds := range favoriteGameIds {
		if len(gameIds) == 0 {
			topGames, err := topGames(DEFAULT_GAME_LIMIT)
			if err != nil {
				fmt.Println("Error occurs in topGames() called by recommendItemsByUserHistory()")
				return nil, err
			}
			recommendedItemMap[itemType], err = recommendByTopGames(itemType, topGames)
			if err != nil {
				fmt.Println(fmt.Println("Error occurs in recommendByTopGames() called by recommendItemsByUserHistory()"))
				return nil, err
			}
		} else {
			recommendedItemMap[itemType], err = recommendByFavoriteHistory(itemType, favoriteItemIds, favoriteGameIds[itemType])
			if err != nil {
				fmt.Println(fmt.Println("Error occurs in recommendByFavoriteHistory() called by recommendItemsByUserHistory()"))
				return nil, err
			}
		}
	}
	return recommendedItemMap, nil
}

// Return a map of Item objects for the given item type as key(one of [Stream, Video, Clip]).
// Each key is corresponding to a list of Items objects, each item object is a recommended item
// based on the top games currently on Twitch.
func (rh RecommendHandler) recommendItemsByDefault() (map[string][]models.Item, error) {
	recommendedItemMap := make(map[string][]models.Item)
	for _, itemType := range models.ItemType {
		recommendedItemMap[itemType] = []models.Item{}
	}

	topGames, err := topGames(DEFAULT_GAME_LIMIT)
	if err != nil {
		fmt.Println("Error occurs in topGames() called by recommendItemsByDefault()")
		return nil, err
	}

	for _, itemType := range models.ItemType {
		recommendedItemMap[itemType], err = recommendByTopGames(itemType, topGames)
		if err != nil {
			fmt.Println("Error occurs in recommendByTopGames() called by recommendItemsByDefault()")
			return nil, err
		}
	}
	return recommendedItemMap, nil
}
