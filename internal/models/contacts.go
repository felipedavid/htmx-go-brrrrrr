package models

import "database/sql"

type Contact struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type ContactModel struct {
	DB *sql.DB
}

func (m ContactModel) GetAll() ([]Contact, error) {
	query := `
		SELECT
			id,
			first_name,
			last_name,
			phone,
			email
		FROM contacts
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	contacts := make([]Contact, 0, 1024)

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

func (m ContactModel) Search(searchTerm string) ([]Contact, error) {
	query := `
		SELECT
			id,
			first_name,
			last_name,
			phone,
			email
		FROM contacts
		WHERE first_name LIKE '%' || $1 || '%'
		OR last_name LIKE '%' || $1 || '%'
		OR phone LIKE '%' || $1 || '%'
		OR email LIKE '%' || $1 || '%'
	`

	rows, err := m.DB.Query(query, searchTerm)
	if err != nil {
		return nil, err
	}

	contacts := make([]Contact, 0, 1024)

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
