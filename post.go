package main

import (
	"mime/multipart"
	"reflect"

	"github.com/olivere/elastic/v7"
)

const (
	POST_INDEX = "post"
)

type Post struct {
	// More info at https://golang.org/pkg/encoding/json/
	// Taking Id as an example, field appears in JSON as key "id".
	// The back-tick syntax frees you from escaping special chars.
	// The current usage is called Tags in Go. Tags can only be used for string values only and can be implemented by using back-ticks
	// Example of using tags to assign default values to strings: https://www.geeksforgeeks.org/how-to-assign-default-value-for-struct-field-in-golang/
	Id      string `json:"id"`
	User    string `json:"user"`
	Message string `json:"message"`
	Url     string `json:"url"`
	Type    string `json:"type"`
}

func searchPostsByUser(user string) ([]Post, error) {
	query := elastic.NewTermQuery("user", user)
	searchResult, err := readFromES(query, POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func searchPostsByKeywords(keywords string) ([]Post, error) {
	query := elastic.NewMatchQuery("message", keywords)
	query.Operator("AND")
	if keywords == "" {
		query.ZeroTermsQuery("all")
	}
	searchResult, err := readFromES(query, POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []Post {
	var ptype Post
	var posts []Post

	// more examples can be found here: https://github.com/olivere/elastic/blob/release-branch.v7/example_test.go
	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		p := item.(Post)
		posts = append(posts, p)
	}
	return posts
}

func savePost(post *Post, file multipart.File) error {
	medialink, err := saveToGCS(file, post.Id)
	if err != nil {
		return err
	}
	post.Url = medialink

	// TODO: Here we may want to improve with a roll back: if the saveToES failed, we want to delete the media saved to GCS above.
	return saveToES(post, POST_INDEX, post.Id)
}

func deletePost(id string, user string) error {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("id", id))
	query.Must(elastic.NewTermQuery("user", user))

	return deleteFromES(query, POST_INDEX)
}
