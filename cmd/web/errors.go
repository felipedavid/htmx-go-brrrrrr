package main

import (
	"log/slog"
	"net/http"
)

func (app *application) ServerError(w http.ResponseWriter, err error) {
	slog.Error("Internal Server Error", "err", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
