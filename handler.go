package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parsing the request body for a json object, schema defined in Post struct.
	fmt.Println("Received one post request")

	// Access-Control-Allow-Headers： used in response to a preflight request to indicate which HTTP headers can be used during the actual request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	// Early return to reject "OPTIONS" method, as both POST and OPTIONS are allowed for this API.
	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var p Post
	if err := decoder.Decode(&p); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "Post received: %s\n", p.Message)
}
