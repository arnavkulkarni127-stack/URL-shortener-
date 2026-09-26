package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error  string `json: "error"`
	Status int    `json: "status"`
}

func WriteErrors(w http.ResponseWriter, err string, status int) {
	response := ErrorResponse{
		Error:  err,
		Status: status,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
