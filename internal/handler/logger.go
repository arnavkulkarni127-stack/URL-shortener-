package handler

import (
	"encoding/json"
	"log"
	"time"
)

type Logs struct {
	Timestamp  string `json: "timestamp"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Status     int    `json:"status"`
	DurationMS int64  `json:"duration_ms"`
}

func WriteLogs(method, path string, status int, duration time.Duration) {
	entry := Logs{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Method:     method,
		Path:       path,
		Status:     status,
		DurationMS: int64(duration),
	}
	jsonData, err := json.Marshal(entry) // returns []byte containing data
	if err != nil {
		return
	}
	log.Println(string(jsonData))
}
