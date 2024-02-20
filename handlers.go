package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusPermanentRedirect)
}

type Contact struct {
	ID        int
	FirstName string
	LastName  string
	Phone     string
	Email     string
}

func listContacts(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("q")

	rows, err := db.Query(`
		SELECT
			id,
			fname,
			lname,
			phone,
			email
		FROM contacts
	`, search)
	handleServerError(w, err)

	var contacts []Contact
	var contact Contact
	for rows.Next() {
		err = rows.Scan(
			&contact.ID,
			&contact.FirstName,
			&contact.LastName,
			&contact.Phone,
			&contact.Email,
		)
		handleServerError(w, err)

		contacts = append(contacts, contact)
	}

	err = rows.Err()
	handleServerError(w, err)

	fmt.Fprintf(w, "%+v", contacts)
}

func handleServerError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		panic(err)
	}
}
