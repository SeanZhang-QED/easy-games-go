package main

import (
	"fmt"
	"log"
    "net/http"

	"github.com/gorilla/mux"
)

func main() {
    fmt.Println("started-service")

	r := mux.NewRouter()

    r.Handle("/game", http.HandlerFunc(getGameHandler)).Methods("Get")

    log.Fatal(http.ListenAndServe(":5000", r))
}