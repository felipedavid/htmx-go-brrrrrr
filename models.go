package main

import (
	"fmt"
	"strings"
)

type Contact struct {
	ID        int
	FirstName string
	LastName  string
	Phone     string
	Email     string

	ValidationErrors map[string]string
}

func NotEmpty(name, value string, errors map[string]string) {
	if len(value) == 0 {
		errors[name] = "cannot be empty"
	}
}

func (c *Contact) Valid() bool {
	errors := c.ValidationErrors

	NotEmpty("FirstName", c.FirstName, errors)
	NotEmpty("LastName", c.LastName, errors)
	NotEmpty("Phone", c.Phone, errors)
	NotEmpty("Email", c.Email, errors)

	return len(errors) == 0
}

func getContact(id int) (Contact, error) {
	query := `
		SELECT id, fname, lname, phone, email from contacts WHERE id = $1
	`

	row := db.QueryRow(query, id)

	contact := Contact{}
	err := row.Scan(
		&contact.ID,
		&contact.FirstName,
		&contact.LastName,
		&contact.Phone,
		&contact.Email,
	)
	if err != nil {
		return contact, err
	}

	return contact, nil
}

func getContacts(searchTerm string) ([]Contact, error) {
	query := `
		SELECT id, fname, lname, phone, email from contacts WHERE
			LOWER(fname) LIKE $1 OR
			LOWER(lname) LIKE $1 OR
			LOWER(phone) LIKE $1 OR
			LOWER(email) LIKE $1
	`

	rows, err := db.Query(query, fmt.Sprintf("%%%s%%", strings.ToLower(searchTerm)))
	if err != nil {
		return nil, err
	}

	var contacts []Contact
	for rows.Next() {
		contact := Contact{}
		err = rows.Scan(
			&contact.ID,
			&contact.FirstName,
			&contact.LastName,
			&contact.Phone,
			&contact.Email,
		)
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func updateContact(contact *Contact) error {
	query := `
		UPDATE contacts SET fname = $1, lname = $2, phone = $3, email = $4 WHERE id = $5
	`

	_, err := db.Exec(query, contact.FirstName, contact.LastName, contact.Phone, contact.Email, contact.ID)
	if err != nil {
		return err
	}

	return nil
}

func deleteContact(id int) error {
	query := `
		DELETE FROM contacts WHERE id = $1
	`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func insertContact(contact *Contact) error {
	query := `
		INSERT INTO contacts (fname, lname, phone, email) VALUES ($1, $2, $3, $4)
	`

	_, err := db.Exec(query, contact.FirstName, contact.LastName, contact.Phone, contact.Email)
	if err != nil {
		return err
	}

	return nil
}
