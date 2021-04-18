package main

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
