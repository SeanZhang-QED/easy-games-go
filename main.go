package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SeanZhang-QED/easy-games-go/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

func main() {
	fmt.Println("started-service")

	mgoSession := getSession()
	uh := handlers.NewUserHandler(mgoSession)
	fh := handlers.NewFavoriteHandler(mgoSession)
	rh := handlers.NewRecommendHandler(mgoSession)

	r := mux.NewRouter()

	r.Handle("/game", http.HandlerFunc(handlers.GetGameHandler)).Methods("GET", "OPTIONS")
	r.Handle("/search", http.HandlerFunc(handlers.SearchItemByGameHandler)).Methods("GET", "OPTIONS")
	r.Handle("/signup", http.HandlerFunc(uh.SignUp)).Methods("POST", "OPTIONS")
	r.Handle("/login", http.HandlerFunc(uh.Login)).Methods("POST", "OPTIONS")
	r.Handle("/logout", http.HandlerFunc(uh.Logout)).Methods("POST", "OPTIONS")
	r.Handle("/favorite", http.HandlerFunc(fh.SetFavorite)).Methods("POST", "OPTIONS")
	r.Handle("/favorite", http.HandlerFunc(fh.UnsetFavorite)).Methods("DELETE", "OPTIONS")
	r.Handle("/favorite", http.HandlerFunc(fh.GetFavorite)).Methods("GET", "OPTIONS")
	r.Handle("/recommendation", http.HandlerFunc(rh.Recommendation)).Methods("GET", "OPTIONS")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}
