package main

import "github.com/felipedavid/contacts/internal/models"

type templateData struct {
	Contacts []models.Contact
	Contact  models.Contact
	Errors   map[string]string
	Flash    string
}
