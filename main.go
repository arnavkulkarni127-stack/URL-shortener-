package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var store *URLStore

func main() {
	conStr := "postgres://postgres:1207@localhost:5432/urlshortener?sslmode=disable" // a connection string to connect to the database

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
	store = &URLStore{db: db}                // initialize the store with the database connection

	//handler: a function that go calls when a request is made to the server
	// that function takes two parameters: a ResponseWriter and a Request
	http.HandleFunc("POST /shorten", shortenHandler) // route for shortening URLs
	http.HandleFunc("GET /{code}", redirectHandler)  // route for redirecting to the original URL || also {code} is a wildcard that will hold annything

	http.ListenAndServe(":8080", nil)
}
