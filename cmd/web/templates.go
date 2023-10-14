package main

import "github.com/felipedavid/contacts/internal/models"

type templateData struct {
	Contacts []models.Contact
	Contact  models.Contact
	Search   string
	Errors   map[string]string
	Flash    string
}
