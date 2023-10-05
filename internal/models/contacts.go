package models

import (
	"database/sql"

	"github.com/felipedavid/contacts/internal/validator"
)

type Contact struct {
	ID        int     `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
	Email     *string `json:"email"`
}

func ValidateContact(v *validator.Validator, contact *Contact) {
	v.Check(contact.FirstName != nil, "first_name", "First Name is required")
	v.Check(contact.LastName != nil, "last_name", "Last Name is required")
	v.Check(contact.Email != nil, "email", "Email is required")
	v.Check(contact.Phone != nil, "phone", "Phone is required")
}

type ContactModel struct {
	DB *sql.DB
}

func (m ContactModel) Get(id int) (*Contact, error) {
	query := `
		SELECT
			id,
			first_name,
			last_name,
			phone,
			email
		FROM contacts
		WHERE id = $1
	`

	row := m.DB.QueryRow(query, id)

	contact := Contact{}

	err := row.Scan(
		&contact.ID,
		&contact.FirstName,
		&contact.LastName,
		&contact.Phone,
		&contact.Email,
	)

	if err != nil {
		return nil, err
	}

	return &contact, nil
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

func (m ContactModel) Insert(contact *Contact) error {
	query := `
		INSERT INTO contacts (
			first_name,
			last_name,
			phone,
			email
		) VALUES (
			$1, $2, $3, $4
		)
	`

	_, err := m.DB.Exec(query,
		contact.FirstName,
		contact.LastName,
		contact.Phone,
		contact.Email,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m ContactModel) Update(contact *Contact) error {
	query := `
		UPDATE contacts
		SET
			first_name = $1,
			last_name = $2,
			phone = $3,
			email = $4
		WHERE id = $5
	`

	_, err := m.DB.Exec(query,
		contact.FirstName,
		contact.LastName,
		contact.Phone,
		contact.Email,
		contact.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m ContactModel) Delete(id int) error {
	query := `
		DELETE FROM contacts
		WHERE id = $1
	`

	_, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
