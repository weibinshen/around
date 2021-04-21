package main

import (
	"fmt"
	"log"
	"net/http"

	// https://github.com/gorilla/mux#examples
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("started-service")

	r := mux.NewRouter()

	// "OPTIONS" is allowed for preflight request for CORS
	r.Handle("/upload", http.HandlerFunc(uploadHandler)).Methods("POST", "OPTIONS")
	r.Handle("/search", http.HandlerFunc(searchHandler)).Methods("GET", "OPTIONS")

	log.Fatal(http.ListenAndServe(":8080", r))
}
