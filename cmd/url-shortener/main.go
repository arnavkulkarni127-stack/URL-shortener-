package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/arnavkulkarni127-stack/URL-shortener-/internal/handler"
	"github.com/arnavkulkarni127-stack/URL-shortener-/internal/store"
	_ "github.com/lib/pq"
)

func main() {
	secret := "my-super-secret-key"
	conStr := "postgres://postgres:1207@localhost:5432/url_shortener?sslmode=disable" // a connection string to connect to the database

	db, err := sql.Open("postgres", conStr) // open a connection pool to the database
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping() // check if the connection is successful

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()                         // close the pool when the main function exits
	log.Println("Connected to the database") // log that the connection was successful

	// initialize the store with the database connection
	s := store.NewURLStore(db)         // Create store
	h := handler.NewHandler(s, secret) // Inject store into handler

	//handler: a function that go calls when a request is made to the server
	// that function takes two parameters: a ResponseWriter and a Request
	http.HandleFunc("POST /api/v1/shorten", h.ShortenHandler)         // route for shortening URLs
	http.HandleFunc("GET /api/v1/redirect/{code}", h.RedirectHandler) // route for redirecting to the original URL || also {code} is a wildcard that will hold annything
	http.HandleFunc("POST /api/v1/auth/signup", h.SignUpHandler)
	http.HandleFunc("POST /api/v1/auth/login", h.LoginHandler)
	// API: set of endpoints(routes) your server exposes to the client to use

	http.ListenAndServe(":8080", nil)
}
