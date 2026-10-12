package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/arnavkulkarni127-stack/URL-shortener-/internal/auth"
)

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc { //take the next hanndler function and wraps it into a middleware 
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization") // get the authorization header from req
		if authHeader == "" {
			WriteErrors(w, "Missing authorization error", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			WriteErrors(w, "Missing Authorization error", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		userID, err := auth.ValidateJWT(tokenString, h.secret)
		if err != nil {
			WriteErrors(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		r.Header.Set("X-User-ID", strconv.Itoa(userID)) // modify the req haeder to a custom hheader and store the userid as string in it.
		next(w, r)

	}

}
