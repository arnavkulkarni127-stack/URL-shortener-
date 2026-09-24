package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	start := time.Now()
	var statusCode int
	var req models.ShortenRequest

	err := json.NewDecoder(r.Body).Decode(&req) // decode the request body into the request struct
	if err != nil {
		statusCode = http.StatusBadRequest
		WriteErrors(w, "Invalid request body", statusCode)

		WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
		return
	}
	if !validation.IsValidURL(req.URL) {
		statusCode = http.StatusBadRequest
		WriteErrors(w, "Invalid URL", statusCode)

		WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
		return
	}
	generatedCode, err := h.store.ShortenRequest(req.URL) // generate a unique short code for the URL
	if err != nil {
		statusCode = http.StatusInternalServerError
		WriteErrors(w, "Failed to generate short code", statusCode)

		WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
		return
	}
	response := models.ShortenResponse{
		ShortCode: generatedCode,
		ShortURL:  fmt.Sprintf("http://localhost:8080/%s", generatedCode), // create the short URL using the generated code
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	statusCode = http.StatusOK
	WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))

}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) { // work is to redirect the user to the original URL based on the code in the path
	start := time.Now()
	var statusCode int
	code := r.PathValue("code")
	url, exists := h.store.GET(code) // get the original URL from the store based on the code
	if !exists {
		statusCode := http.StatusNotFound
		WriteErrors(w, "Short code not found", statusCode)
		WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
		return
	}

	statusCode = http.StatusFound
	WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
	http.Redirect(w, r, url, statusCode) // write the status code to the response, which tells the browser that the resource has been found and is being redirected
}
