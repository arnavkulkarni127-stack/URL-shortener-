package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arnavkulkarni127-stack/URL-shortener-/internal/models"
	"github.com/arnavkulkarni127-stack/URL-shortener-/internal/store"
	"github.com/arnavkulkarni127-stack/URL-shortener-/pkg/validation"
)

type URLStore struct {
	// mutex to protect the map from concurrent access (like a lock)
	Db *sql.DB
}
type Handler struct {
	store *store.URLStore
}

func NewHandler(s *store.URLStore) *Handler {
	return &Handler{store: s}
}
func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {

	err := json.NewDecoder(r.Body).Decode(&models.Request) // decode the request body into the request struct
	if err != nil {
		WriteErrors(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if !validation.IsValidURL(models.Request.URL) {
		WriteErrors(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	generatedCode, err := h.store.ShortenRequest(models.Request.URL) // generate a unique short code for the URL
	if err != nil {
		WriteErrors(w, "Failed to generate short code", http.StatusInternalServerError)
		return
	}
	response := models.ShortenResponse{
		ShortCode: generatedCode,
		ShortURL:  fmt.Sprintf("http://localhost:8080/%s", generatedCode), // create the short URL using the generated code
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) { // work is to redirect the user to the original URL based on the code in the path
	code := r.PathValue("code")
	url, exists := h.store.GET(code) // get the original URL from the store based on the code
	if !exists {
		WriteErrors(w, "Short code not found", http.StatusNotFound)
		return
	}

	fmt.Printf("Redirecting code: %s\n", code)
	http.Redirect(w, r, url, http.StatusFound) // write the status code to the response, which tells the browser that the resource has been found and is being redirected
}
