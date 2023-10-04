package main

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, err := template.New(page).ParseFiles("./ui/html/base.tmpl")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	ts, err = ts.ParseFiles(fmt.Sprintf("./ui/html/pages/%s.tmpl", page))
	if err != nil {
		app.ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}
