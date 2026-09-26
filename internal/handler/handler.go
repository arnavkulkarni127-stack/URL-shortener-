package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arnavkulkarni127-stack/URL-shortener-/internal/auth"
	"github.com/arnavkulkarni127-stack/URL-shortener-/internal/models"
	"github.com/arnavkulkarni127-stack/URL-shortener-/internal/store"
	"github.com/arnavkulkarni127-stack/URL-shortener-/pkg/validation"
)

type URLStore struct {
	// mutex to protect the map from concurrent access (like a lock)
	Db *sql.DB
}
type Handler struct {
	store  *store.URLStore
	secret string
}

func NewHandler(s *store.URLStore, secret string) *Handler {
	return &Handler{store: s, secret: secret}
}
func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var statusCode int
	var req models.ShortenRequest
	//var userCreated models.CreateUser

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
	generatedCode, err := h.store.ShortenRequest(1, req.URL) // generate a unique short code for the URL
	if err != nil {
		statusCode = http.StatusInternalServerError
		fmt.Println("ShortenRequest error:", err)
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

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var userCreated struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	start := time.Now()
	//	Decode the request
	err := json.NewDecoder(r.Body).Decode(&userCreated)
	if err != nil {
		statusCode = http.StatusBadRequest
		WriteErrors(w, "Invalid request body", statusCode)

		WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
		return
	}
	//  password hash
	passHash, err := auth.HashPassword(userCreated.Password)
	if err != nil {
		statusCode = http.StatusInternalServerError
		WriteErrors(w, "failed to hash passwords", statusCode)
		WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
		return

	}
	userID, err := h.store.CreateUser(userCreated.Email, passHash)
	if err != nil {
		statusCode = http.StatusBadRequest
		WriteErrors(w, "Email already exists", statusCode)
		WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
		return
	}

	tokenString, erro := auth.GenerateJWT(userID, h.secret)
	if erro != nil {
		statusCode = http.StatusInternalServerError
		WriteErrors(w, "Failed to generate token", statusCode)
		WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
		return
	}

	statusCode = http.StatusCreated
	response := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: tokenString,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
}

func (h Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userCreated struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var statusCode int
	start := time.Now()
	//	Decode the request
	err := json.NewDecoder(r.Body).Decode(&userCreated)
	if err != nil {
		statusCode = http.StatusBadRequest
		WriteErrors(w, "Invalid request body", statusCode)

		WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
		return
	}
	//	get the user by email
	_, passHash, err := h.store.GetUserByEmail(userCreated.Email)
	if err != nil {
		statusCode = http.StatusUnauthorized
		WriteErrors(w, "Invalid Credentials", statusCode)
		WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
	}
	// 	verify the password
	if !auth.VerifyPassword(passHash, userCreated.Password) {
		statusCode = http.StatusUnauthorized
		WriteErrors(w, "Invalid credentials", statusCode)
		WriteLogs(r.Method, r.RequestURI, statusCode, time.Since(start))
	}

}
