package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func serverErrorResponse(w http.ResponseWriter) {
	_, _ = w.Write([]byte("shit hit the fan"))
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
	mux.HandleFunc("POST /contacts/new", addContactHandler)

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

func addContactHandler(w http.ResponseWriter, r *http.Request) {
	var contact Contact

	contact.FirstName = r.PostForm.Get("FirstName")
	contact.LastName = r.PostForm.Get("LastName")
	contact.Phone = r.PostForm.Get("Phone")
	contact.Email = r.PostForm.Get("email")

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

	err := renderPage(w, "contact", map[string]any{
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
