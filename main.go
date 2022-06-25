package main

import (
	"fmt"
	"log"
    "net/http"

	"github.com/gorilla/mux"
	"github.com/SeanZhang-QED/easy-games-go/handlers"
)

func main() {
    fmt.Println("started-service")

	r := mux.NewRouter()

    r.Handle("/game", http.HandlerFunc(handlers.GetGameHandler)).Methods("GET", "OPTIONS")

    log.Fatal(http.ListenAndServe(":8080", r))
}