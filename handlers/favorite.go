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
	err = fh.checkItem(favorite.Item)	
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
	fmt.Printf("Received item: %v\n", favorite)
	err = fh.unsetFavoriteItem(loggedSession.Email, favorite.Item.Id)
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
	items, err := fh.getItemByType(loggedSession.Email)
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

func (fh FavoriteHandler) checkItem(item models.Item) error {
	var items []models.Item
	err := fh.session.DB("easy-games-db").C("items").Find(bson.M{"_id": item.Id}).All(&items)
	if err != nil {
		fmt.Println("Fail to search the item info from Items document")
		return err
	}

	if len(items) == 0 {
		err = fh.session.DB("easy-games-db").C("items").Insert(item)
		if err != nil {
			fmt.Println("Fail to add the item into Items document")
			return err
		}
		return nil
	}
	return nil
}

func (fh FavoriteHandler) setFavoriteItem(email string, itemId string) error {
	var u models.User
	err := fh.session.DB("easy-games-db").C("users").FindId(email).One(&u)
	if err != nil {
		fmt.Println("Fail to fetch the user info from users document")
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

func (fh FavoriteHandler) unsetFavoriteItem(email string, itemId string) error {
	var u models.User
	err := fh.session.DB("easy-games-db").C("users").FindId(email).One(&u)
	if err != nil {
		fmt.Println("Fail to fetch the user info from users document")
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

func (fh FavoriteHandler) getItemByType(email string) (map[string][]models.Item, error) {
	m := make(map[string][]models.Item)

	// get the user 
	var u models.User
	err := fh.session.DB("easy-games-db").C("users").FindId(email).One(&u)
	if err != nil {
		fmt.Println("Fail to fetch the user info from users document")
		return m, err
	}
	// iterate the favorite recods map
	for itemId := range u.FavoriteRecords {
		item, err := fh.getItemByItemId(itemId)
		if err != nil {
			fmt.Println("Fail to fetch the item info by id")
			return m, err
		}
		m[item.ItemType] = append(m[item.ItemType], item)
	}

	return m, nil
}

func (fh FavoriteHandler) getItemByItemId(itemId string) (models.Item, error) {
	var curItem models.Item
	err := fh.session.DB("easy-games-db").C("items").Find(bson.M{"_id": itemId}).One(&curItem)
	if err != nil {
		fmt.Println("Fail to fetch the item info by id")
		return models.Item{}, err
	}
	return curItem, nil
}