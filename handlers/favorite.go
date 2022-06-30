package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SeanZhang-QED/easy-games-go/models"
	"github.com/SeanZhang-QED/easy-games-go/session"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// added session to our userController
type FavoriteHandler struct {
	session *mgo.Session
}

// added session to our userController
func NewFavoriteHandler(s *mgo.Session) *FavoriteHandler {
	return &FavoriteHandler{s}
}

func (fh FavoriteHandler) SetFavorite(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one Set Favorite Item request")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	// step 1: check for login
	_, err := r.Cookie("sessionId")
	if err != nil {
		http.Error(w, "Haven't logged in", http.StatusUnauthorized)
		fmt.Println("Fail to read cookie from http request")
		return 
	}
	loggedSession, err := session.AlreadyLoggedIn(w, r, fh.session)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	} 
	// step 2: unmarshal the favorite item object from JSON
	decoder := json.NewDecoder(r.Body)
	var favorite models.Favorite
	if err := decoder.Decode(&favorite); err != nil {
		http.Error(w, "Cannot decode favorite Item from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode favorite Item from client %v\n", err)
		return
	}
	// fmt.Printf("Received item: %v\n", favorite)
	// step 3: insert/get the foavoirte item from the mongoDB
	err = checkItem(favorite.Item, fh.session)	
	if err != nil {
		http.Error(w, "Fail to check the item info in MongoDB", http.StatusInternalServerError)
		fmt.Printf("Fail to check the item info in MongoDB: %v\n", err)
		return
	}
	err = fh.setFavoriteItem(loggedSession.Email, favorite.Item.Id)
	if err != nil {
		http.Error(w, "Fail to add the item info into User's favorite list", http.StatusInternalServerError)
		fmt.Printf("Fail to add the item info into User's favorite list: %v\n", err)
		return
	}
	fmt.Printf("Add item as favorite item successfully: %s.\n", loggedSession.Email)
}

func (fh FavoriteHandler) UnsetFavorite(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one Delte Favorite Item request")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/plain")

	if r.Method == http.MethodOptions {
		return
	}

	// step 1: check for login
	_, err := r.Cookie("sessionId")
	if err != nil {
		http.Error(w, "Haven't logged in", http.StatusUnauthorized)
		fmt.Println("Fail to read cookie from http request")
		return 
	}
	loggedSession, err := session.AlreadyLoggedIn(w, r, fh.session)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	} 
	// step 2: unmarshal the favorite item object from JSON
	decoder := json.NewDecoder(r.Body)
	var favorite models.Favorite
	if err := decoder.Decode(&favorite); err != nil {
		http.Error(w, "Cannot decode favorite Item from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode favorite Item from client %v\n", err)
		return
	}
	err = fh.deleteFavoriteItem(loggedSession.Email, favorite.Item.Id)
	if err != nil {
		http.Error(w, "Fail to add the item info into User's favorite list", http.StatusInternalServerError)
		fmt.Printf("Fail to add the item info into User's favorite list: %v\n", err)
		return
	}
	fmt.Printf("Delete item from favorite item successfully: %s.\n", loggedSession.Email)
}

func (fh FavoriteHandler) GetFavorite(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one Get Favorite Item List request")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	// step 1: check for login
	_, err := r.Cookie("sessionId")
	if err != nil {
		http.Error(w, "Haven't logged in", http.StatusUnauthorized)
		fmt.Println("Fail to read cookie from http request")
		return 
	}
	loggedSession, err := session.AlreadyLoggedIn(w, r, fh.session)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	// step 2: get a map of items, key is type, value is a list of Item
	items, err := fh.getFavoriteItemsByType(loggedSession.Email)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		fmt.Printf("Fail to get the favorite items: %v", err)
		return
	}

	// step 3: return the map as response json body
	// Marshal: Game -> JSON
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		// return msg with http status code
		http.Error(w, "Failed to convert items map to JSON", http.StatusInternalServerError)
		fmt.Printf("Failed to convert items map to JSON: %v. \n", err)
		return
	}
	w.Write(itemsJSON)
	fmt.Printf("Get user's all favorite items successfully: %s.\n", loggedSession.Email) 
}

func (fh FavoriteHandler) setFavoriteItem(email string, itemId string) error {
	u, err := getUserByEmail(email, fh.session)
	if ; err != nil {
		return err
	}

	if _ , ok := u.FavoriteRecords[itemId]; ok {
		fmt.Println("Already add the item as Favorite Item")
		return nil
	} else {
		var value models.Void
		u.FavoriteRecords[itemId] = value
		err = fh.session.DB("easy-games-db").C("users").UpdateId(email, bson.M{"$set": bson.M{"favorite_records": u.FavoriteRecords}})
		if err != nil {
			fmt.Println("Fail to update the user info to MongoDB")
			return err
		}
	}
	return nil
}

func (fh FavoriteHandler) deleteFavoriteItem(email string, itemId string) error {
	u, err := getUserByEmail(email, fh.session)
	if ; err != nil {
		return err
	}
	if _ , ok := u.FavoriteRecords[itemId]; ok {
		delete(u.FavoriteRecords, itemId)
		err = fh.session.DB("easy-games-db").C("users").UpdateId(email, bson.M{"$set": bson.M{"favorite_records": u.FavoriteRecords}})
		if err != nil {
			fmt.Println("Fail to update the user info to MongoDB")
			return err
		}
	} else {
		fmt.Println("Item is not in the Favorite Item records")
		return nil
	}	
	return nil
}

func (fh FavoriteHandler) getFavoriteItemsByType(email string) (map[string][]models.Item, error) {
	m := make(map[string][]models.Item)

	// get the user 
	u, err := getUserByEmail(email, fh.session)
	if ; err != nil {
		return nil, err
	}
	// iterate the favorite recods map
	for itemId := range u.FavoriteRecords {
		item, err := getItemByItemId(itemId, fh.session)
		if err != nil {
			fmt.Println("Fail to fetch the item info by id")
			return m, err
		}
		m[item.ItemType] = append(m[item.ItemType], item)
	}

	return m, nil
}

// Get the set of favoriteItem ids for given User
func getFavoriteItemIds(email string, session *mgo.Session) (map[string]models.Void, error) {
	u, err := getUserByEmail(email, session)
	if err != nil {
		fmt.Println("Fail to get the User by email.")
		return nil, err
	}
	return u.FavoriteRecords, nil
}

// Get favoriteItems's Game id for the given user.
func getFavoriteGameIds(email string, session *mgo.Session) (map[string][]string, error) {
	
	// Step 1: get the user's favorite Item ids
	itemIds, err := getFavoriteItemIds(email, session)
	if err != nil {
		fmt.Println("Fail to get the User's Favorite Items info from DB")
		return nil, err
	}
	// Step 2: prepare the map
	m := make(map[string][]string)
	for _, itemType := range models.ItemType {
		m[itemType] = []string{}
	}
	// Step 3: iterate the set and allocate the map(ItemType To a slice GameId)
	for itemId := range itemIds {
		curItem, err := getItemByItemId(itemId, session)
		if err != nil {
			return nil, err
		}
		m[curItem.ItemType] = append(m[curItem.ItemType], curItem.GameId)
	}
	return m, nil
}