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
	// I put "OPTIONS" here as an example of how to handle different HTTP methods to the same API in uploadHandler.
	r.Handle("/upload", http.HandlerFunc(uploadHandler)).Methods("POST", "OPTIONS")

	log.Fatal(http.ListenAndServe(":8080", r))
}
