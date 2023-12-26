package main

import (
	"database/sql"
	"net/http"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	conn, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		panic(err)
	}

	if err = conn.Ping(); err != nil {
		panic(err)
	}

	_, _ = conn.Exec(`
		CREATE TABLE contacts (
		    id INTEGER PRIMARY KEY,
			fname TEXT NOT NULL,
			lname TEXT NOT NULL,
			phone TEXT NOT NULL,
			email TEXT NOT NULL
		);
	`)

	db = conn
}

func main() {
	mux := defineRoutes()

	err := http.ListenAndServe(":8080", mux)

	panic(err)
}
