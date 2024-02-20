package main

import (
	"database/sql"
	"net/http"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {
	connectToDatabase()
	applySchema()
	startServer()
}

func connectToDatabase() {
	var err error
	db, err = sql.Open("sqlite", "sqlite.db")
	shouldNotErrorOut(err)
}

func applySchema() {
	schema := `
		CREATE TABLE IF NOT EXISTS contacts (
		    id INTEGER PRIMARY KEY,
			fname TEXT NOT NULL,
			lname TEXT NOT NULL,
			phone TEXT NOT NULL,
			email TEXT NOT NULL                                    
		) STRICT;
	`

	_, err := db.Exec(schema)
	shouldNotErrorOut(err)
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			recover()
		}()

		next.ServeHTTP(w, r)
	})
}

func startServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", index)
	mux.HandleFunc("GET /contacts", listContacts)

	err := http.ListenAndServe("127.0.0.1:6969", recoverPanic(mux))
	shouldNotErrorOut(err)
}

func shouldNotErrorOut(err error) {
	if err != nil {
		panic(err)
	}
}
