package store

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/arnavkulkarni127-stack/URL-shortener-/pkg/codegen"
)

type URLStore struct {
	db *sql.DB
}

func NewURLStore(db *sql.DB) *URLStore {
	return &URLStore{db: db}
}

func (store *URLStore) GET(code string) (string, bool) { // get the long url from database

	var OriginalURL string
	err := store.db.QueryRow("SELECT original_url FROM urls WHERE code = $1", code).Scan(&OriginalURL) // query the database for the original URL corresponding to the given code

	if err == sql.ErrNoRows {
		return "", false
	}
	if err != nil {
		fmt.Println("Error retrieving original URL:", err)
		return "", false
	}

	return OriginalURL, true
}

func (store *URLStore) ShortenRequest(ogurl string) (string, error) {
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		code := codegen.GenerateRandomShortCode(6)                                                   // generate a random short code of length 6
		_, err := store.db.Exec("INSERT INTO urls(code, original_url) VALUES ($1, $2)", code, ogurl) // insert query executed

		if err == nil {
			return code, nil
		}
		if strings.Contains(err.Error(), "duplicate key") {
			continue
		}

		return "", err
	}

	return "", fmt.Errorf("failed to generate a unique code after %d attempts", maxRetries)
}
