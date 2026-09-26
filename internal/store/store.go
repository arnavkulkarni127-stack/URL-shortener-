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

func (store *URLStore) ShortenRequest(userID int, ogurl string) (string, error) {
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		code := codegen.GenerateRandomShortCode(6)                                                                        // generate a random short code of length 6
		_, err := store.db.Exec("INSERT INTO urls(code, original_url, user_id) VALUES ($1, $2, $3)", code, ogurl, userID) // insert query executed

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
func (store *URLStore) CreateUser(email, passwordHash string) (int, error) {
	var userID int
	err := store.db.QueryRow("INSERT INTO users(email, password_hash) VALUES ($1, $2) RETURNING id", email, passwordHash).Scan(&userID)
	if err != nil {
		fmt.Println("CreateUser error:", err) // Add this
		return 0, err
	}
	return userID, nil
}
func (store *URLStore) GetUserByEmail(email string) (int, string, error) {
	var userID int
	var password string

	err := store.db.QueryRow("SELECT id,password_hash FROM users WHERE email = &1", email).Scan(&userID, &password)
	if err != nil {
		return 0, "", err
	}
	return userID, password, nil

}
