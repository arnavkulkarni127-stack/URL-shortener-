package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func shortenHandler(w http.ResponseWriter, r *http.Request) {

	err := json.NewDecoder(r.Body).Decode(&request) // decode the request body into the request struct
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if !isValidURL(request.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
}
func redirectHandler(w http.ResponseWriter, r *http.Request) { // work is to redirect the user to the original URL based on the code in the path
	code := r.PathValue("code")
	fmt.Printf("Redirecting code: %s\n", code)
	w.Header().Set("Location", "https://github.com") // wriete the location header to the response, which tells the browser where to redirect to
	w.WriteHeader(302)                               // write the status code to the response, which tells the browser that the resource has been found and is being redirected
}
