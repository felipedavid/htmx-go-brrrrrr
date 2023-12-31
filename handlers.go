package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func serverErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("shit hit the fan"))
}

func badRequestResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte("you're dumb or what?"))
}

func renderPage(w http.ResponseWriter, pageName string, data any) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("./templates/pages/%s.gohtml", pageName), "./templates/base.gohtml")
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		return err
	}

	return nil
}

func defineRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", indexPageHandler)
	mux.HandleFunc("GET /contacts", contactsPageHandler)
	mux.HandleFunc("GET /contacts/new", renderAddContactPageHandler)
	mux.HandleFunc("POST /contacts/new", processNewContactForm)
	mux.HandleFunc("GET /contacts/{id}/view", renderContactPageHandler)
	mux.HandleFunc("GET /contacts/{id}/edit", renderEditContactPageHandler)
	mux.HandleFunc("POST /contacts/{id}/edit", processEditContactForm)
	mux.HandleFunc("POST /contacts/{id}/delete", deleteContactHandler)

	return mux
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func contactsPageHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")

	contacts, err := getContacts(searchQuery)
	if err != nil {
		serverErrorResponse(w)
		return
	}

	err = renderPage(w, "index", map[string]any{
		"Contacts":    contacts,
		"SearchQuery": searchQuery,
	})
	if err != nil {
		serverErrorResponse(w)
	}
}

func processNewContactForm(w http.ResponseWriter, r *http.Request) {
	var contact Contact

	err := r.ParseForm()
	if err != nil {
		serverErrorResponse(w)
		return
	}

	contact.FirstName = r.FormValue("FirstName")
	contact.LastName = r.FormValue("LastName")
	contact.Phone = r.FormValue("Phone")
	contact.Email = r.FormValue("Email")

	if !contact.Valid() {
		err := renderPage(w, "new_contact", map[string]any{
			"Contact":          contact,
			"ValidationErrors": contact.ValidationErrors,
		})
		if err != nil {
			serverErrorResponse(w)
			return
		}
	}

	err = insertContact(&contact)
	if err != nil {
		serverErrorResponse(w)
		return
	}

	err = renderPage(w, "contact", map[string]any{
		"Contact": contact,
	})
	if err != nil {
		serverErrorResponse(w)
		return
	}
}

func deleteContactHandler(w http.ResponseWriter, r *http.Request) {
	contactID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequestResponse(w)
		return
	}

	err = deleteContact(contactID)
	if err != nil {
		serverErrorResponse(w)
		return
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func processEditContactForm(w http.ResponseWriter, r *http.Request) {
	contactID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequestResponse(w)
		return
	}

	contact, err := getContact(contactID)
	if err != nil {
		serverErrorResponse(w)
		return
	}

	err = r.ParseForm()
	if err != nil {
		serverErrorResponse(w)
		return
	}

	contact.FirstName = r.FormValue("FirstName")
	contact.LastName = r.FormValue("LastName")
	contact.Phone = r.FormValue("Phone")
	contact.Email = r.FormValue("Email")

	if !contact.Valid() {
		err := renderPage(w, "edit_contact", map[string]any{
			"Contact":          contact,
			"ValidationErrors": contact.ValidationErrors,
		})
		if err != nil {
			serverErrorResponse(w)
			return
		}
	}

	err = updateContact(&contact)
	if err != nil {
		serverErrorResponse(w)
		return
	}

	err = renderPage(w, "contact", map[string]any{
		"Contact": contact,
	})
	if err != nil {
		serverErrorResponse(w)
		return
	}
}

func renderAddContactPageHandler(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "new_contact", map[string]any{
		"Contact": Contact{},
	})
	if err != nil {
		serverErrorResponse(w)
		return
	}
}

func renderContactPageHandler(w http.ResponseWriter, r *http.Request) {
	contactID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequestResponse(w)
		return
	}

	contact, err := getContact(contactID)
	if err != nil {
		serverErrorResponse(w)
		return
	}

	err = renderPage(w, "contact", map[string]any{
		"Contact": contact,
	})
	if err != nil {
		serverErrorResponse(w)
		return
	}
}

func renderEditContactPageHandler(w http.ResponseWriter, r *http.Request) {
	contactID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequestResponse(w)
		return
	}

	contact, err := getContact(contactID)
	if err != nil {
		serverErrorResponse(w)
		return
	}

	err = renderPage(w, "edit_contact", map[string]any{
		"Contact": contact,
	})
	if err != nil {
		serverErrorResponse(w)
		return
	}
}
