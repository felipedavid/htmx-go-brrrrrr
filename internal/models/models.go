package models

import (
	"database/sql"
	"os"
)

type Models struct {
	Contacts ContactModel
}

func New(db *sql.DB) *Models {
	return &Models{
		Contacts: ContactModel{DB: db},
	}
}

func RunMigrations(db *sql.DB) error {
	query, err := os.ReadFile("migrations.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(query))
	if err != nil {
		return err
	}

	return nil
}
