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

	uh := handlers.NewUserHandler(getSession())

	r := mux.NewRouter()

	r.Handle("/game", http.HandlerFunc(handlers.GetGameHandler)).Methods("GET", "OPTIONS")
	r.Handle("/signup", http.HandlerFunc(uh.SignUp)).Methods("POST", "OPTIONS")
	r.Handle("/login", http.HandlerFunc(uh.Login)).Methods("POST", "OPTIONS")

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
