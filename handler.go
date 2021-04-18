package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// An interesting reason by w doesn't have "*", but r does:
// https://stackoverflow.com/questions/13255907/in-go-http-handlers-why-is-the-responsewriter-a-value-but-the-request-a-pointer
// The http.ResponseWriter is an interface, and the existing types implementing this interface are pointers.
// That means there's no need to use a pointer to this interface, as it's already "backed" by a pointer.
// http.Request is not an interface, it's just a struct.
// It is not supposed to be mutable, but it is large so copying it will reduce performance.
// So it is passed in as a pointer.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parsing the request body for a json object, schema defined in Post struct.
	fmt.Println("Received one post request")

	// Access-Control-Allow-Origin: Specify which service is accepted for CORS.
	// Access-Control-Allow-Headers： used in response to a preflight request to indicate which HTTP headers can be used during the actual request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	// Early return to respond to preflight request signified by "OPTIONS" method, with the headers properly set above.
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
