package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.indexView)
	mux.HandleFunc("/contacts", app.contactsHandler)
	mux.HandleFunc("/contacts/new", app.newContactsHandler)
	mux.HandleFunc("/contacts/edit", app.editContactsHandler)
	mux.HandleFunc("/contacts/delete", app.deleteContactsHandler)

	return app.sessionManager.LoadAndSave(mux)
}
