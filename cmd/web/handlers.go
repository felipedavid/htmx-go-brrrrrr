package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/felipedavid/contacts/internal/models"
	"github.com/felipedavid/contacts/internal/validator"
)

func (app *application) indexView(w http.ResponseWriter, r *http.Request) {
	time.Sleep(app.config.artificialDelay)
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (app *application) contactsHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(app.config.artificialDelay)

	queries := r.URL.Query()

	id, err := strconv.Atoi(queries.Get("id"))

	if err == nil && id > 0 {
		contact, err := app.models.Contacts.Get(id)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.render(w, http.StatusOK, "show", &templateData{
			Contact: contact,
		})
		return
	}

	switch r.Method {
	case http.MethodGet:
		searchTerm := queries.Get("q")

		var contacts []models.Contact

		var err error
		if searchTerm == "" {
			contacts, err = app.models.Contacts.GetAll()
		} else {
			contacts, err = app.models.Contacts.Search(searchTerm)
		}

		if err != nil {
			app.ServerError(w, err)
			return
		}

		flashMessage := app.sessionManager.PopString(r.Context(), "flash")
		app.render(w, http.StatusOK, "contacts", &templateData{
			Contacts: contacts,
			Flash:    flashMessage,
			Search:   searchTerm,
		})
	default:
		app.MethodNotAllowed(w)
	}
}

func (app *application) newContactsHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(app.config.artificialDelay)
	switch r.Method {
	case http.MethodGet:
		app.render(w, http.StatusOK, "new", &templateData{})
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			app.ServerError(w, err)
			return
		}

		contact := &models.Contact{}

		contact.FirstName = r.PostFormValue("first_name")
		contact.LastName = r.PostFormValue("last_name")
		contact.Email = r.PostFormValue("email")
		contact.Phone = r.PostFormValue("phone")

		v := validator.New()

		if models.ValidateContact(v, contact); !v.Valid() {
			app.render(w, http.StatusBadRequest, "new", &templateData{
				Contact: contact,
				Errors:  v.Errors,
			})
			return
		}

		if err := app.models.Contacts.Insert(contact); err != nil {
			app.ServerError(w, err)
			return
		}

		app.sessionManager.Put(r.Context(), "flash", "Contact successfully created")
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	default:
		app.MethodNotAllowed(w)
	}
}

func (app *application) editContactsHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(app.config.artificialDelay)
	queries := r.URL.Query()

	id, err := strconv.Atoi(queries.Get("id"))

	if err != nil || id <= 0 {
		app.NotFound(w)
		return
	}

	contact, err := app.models.Contacts.Get(id)

	if err != nil {
		app.ServerError(w, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.render(w, http.StatusOK, "edit", &templateData{
			Contact: contact,
		})
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			app.ServerError(w, err)
			return
		}

		contact.FirstName = r.PostFormValue("first_name")
		contact.LastName = r.PostFormValue("last_name")
		contact.Email = r.PostFormValue("email")
		contact.Phone = r.PostFormValue("phone")

		v := validator.New()

		if models.ValidateContact(v, contact); !v.Valid() {
			app.render(w, http.StatusOK, "edit", &templateData{
				Contact: contact,
				Errors:  v.Errors,
			})
			return
		}

		if err := app.models.Contacts.Update(contact); err != nil {
			app.ServerError(w, err)
			return
		}

		app.sessionManager.Put(r.Context(), "flash", "Contact successfully updated")
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	default:
		app.MethodNotAllowed(w)
	}
}

func (app *application) deleteContactsHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(app.config.artificialDelay)
	if r.Method != http.MethodDelete {
		app.MethodNotAllowed(w)
		return
	}

	queries := r.URL.Query()

	id, err := strconv.Atoi(queries.Get("id"))

	if err != nil || id <= 0 {
		app.NotFound(w)
		return
	}

	if err := app.models.Contacts.Delete(id); err != nil {
		app.ServerError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Contact successfully deleted")
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}
