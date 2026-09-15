package main

import (
	"database/sql"
)

var request struct {
	URL string `json:"url"` // tellsw the decoder to look for a field called url in the request body and store it in the URL field of the request struct
}

type URLStore struct {
	// mutex to protect the map from concurrent access (like a lock)
	db *sql.DB
}
type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}
