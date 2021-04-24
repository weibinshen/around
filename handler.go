package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
)

var (
	mediaTypes = map[string]string{
		".jpeg": "image",
		".jpg":  "image",
		".gif":  "image",
		".png":  "image",
		".mov":  "video",
		".mp4":  "video",
		".avi":  "video",
		".flv":  "video",
		".wmv":  "video",
	}
)

// An interesting reason by w doesn't have "*", but r does:
// https://stackoverflow.com/questions/13255907/in-go-http-handlers-why-is-the-responsewriter-a-value-but-the-request-a-pointer
// The http.ResponseWriter is an interface
// An interface can store either a struct directly or a pointer to a struct,
// and an interface pointer is almost never useful and is discouraged by go.
// https://stackoverflow.com/questions/44370277/type-is-pointer-to-interface-not-interface-confusion
// http.Request is not an interface, it's just a struct.
// It is not supposed to be mutable, but it is large so copying it will reduce performance.
// So it is passed in as a pointer.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parsing the request body for a json object, schema defined in Post struct.
	fmt.Println("Received one upload request")

	// Access-Control-Allow-Origin: Specify which service is accepted for CORS.
	// Access-Control-Allow-Headers： used in response to a preflight request to indicate which HTTP headers can be used during the actual request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	// Early return to respond to preflight request signified by "OPTIONS" method, with the headers properly set above.
	if r.Method == "OPTIONS" {
		return
	}

	p := Post{
		Id:      uuid.New().String(), // Note: Need to convert UUID to String.
		User:    r.FormValue("user"),
		Message: r.FormValue("message"),
	}

	file, header, err := r.FormFile("media_file")
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		fmt.Printf("Media file is not available %v\n", err)
		return
	}

	suffix := filepath.Ext(header.Filename)
	if t, ok := mediaTypes[suffix]; ok {
		p.Type = t
	} else {
		p.Type = "unknown"
	}

	err = savePost(&p, file)
	if err != nil {
		http.Error(w, "Failed to save post to GCS or Elasticsearch", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to GCS or Elasticsearch %v\n", err)
		return
	}

	fmt.Println("Post is saved successfully.")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for search")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords")

	var posts []Post
	var err error
	if user != "" {
		posts, err = searchPostsByUser(user)
	} else {
		posts, err = searchPostsByKeywords(keywords)
	}

	if err != nil {
		http.Error(w, "Failed to read post from Elasticsearch", http.StatusInternalServerError)
		fmt.Printf("Failed to read post from Elasticsearch %v.\n", err)
		return
	}

	js, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
		return
	}
	w.Write(js)
}
