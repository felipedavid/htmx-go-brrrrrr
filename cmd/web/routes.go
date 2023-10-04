package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.indexHandler)
	mux.HandleFunc("/contacts", app.contactsHandler)

	return mux
}
