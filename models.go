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

	return len(errors) > 0
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
