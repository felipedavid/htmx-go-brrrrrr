package main

type Contact struct {
	ID        int
	FirstName string
	LastName  string
	Phone     string
	Email     string
}

func getContacts(searchTerm string) ([]Contact, error) {
	query := `SELECT id, fname, lname, phone, email from contacts`

	rows, err := db.Query(query)
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
		)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}
