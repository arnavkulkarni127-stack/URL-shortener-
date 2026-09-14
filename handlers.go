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
	generatedCode, err := store.ShortenRequest(request.URL) // generate a unique short code for the URL
	if err != nil {
		http.Error(w, "Failed to generate short code", http.StatusInternalServerError)
		return
	}
	response := ShortenResponse{
		ShortCode: generatedCode,
		ShortURL:  fmt.Sprintf("http://localhost:8080/%s", generatedCode), // create the short URL using the generated code
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) { // work is to redirect the user to the original URL based on the code in the path
	code := r.PathValue("code")
	url, exists := store.GET(code) // get the original URL from the store based on the code
	if !exists {
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}

	fmt.Printf("Redirecting code: %s\n", code)
	http.Redirect(w, r, url, http.StatusFound) // write the status code to the response, which tells the browser that the resource has been found and is being redirected
}
