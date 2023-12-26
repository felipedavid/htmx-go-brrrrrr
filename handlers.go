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
	mux.HandleFunc("GET /contacts/{searchTerm}", contactsPageHandler)

	return mux
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "index", nil)
	if err != nil {
		serverErrorResponse(w)
	}
}

func contactsPageHandler(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.PathValue("searchTerm")

	contacts, err := getContacts(searchTerm)
	if err != nil {
		serverErrorResponse(w)
		return
	}

	err = renderPage(w, "contacts", contacts)
	if err != nil {
		serverErrorResponse(w)
		return
	}
}
