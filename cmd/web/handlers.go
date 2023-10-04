package main

import (
	"fmt"
	"net/http"

	"github.com/felipedavid/contacts/internal/models"
)

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (app *application) contactsHandler(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")

	var contacts []models.Contact

	var err error
	if searchTerm == "" {
		contacts, err = app.models.Contacts.GetAll()
	} else {
		contacts, err = app.models.Contacts.Search(searchTerm)
	}

	if err != nil {
		app.ServerError(w)
		return
	}

	fmt.Fprintf(w, "%v", contacts)
}
