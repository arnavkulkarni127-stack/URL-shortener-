package main

import (
	"net/http"
)

func main() {
	//handler: a function that go calls when a request is made to the server
	// that function takes two parameters: a ResponseWriter and a Request
	http.HandleFunc("POST /shorten", shortenHandler) // route for shortening URLs
	http.HandleFunc("GET /{code}", redirectHandler)  // route for redirecting to the original URL || also {code} is a wildcard that will hold annything

	http.ListenAndServe(":8080", nil)
}
