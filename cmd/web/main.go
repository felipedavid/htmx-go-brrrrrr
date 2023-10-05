package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/felipedavid/contacts/internal/models"

	_ "modernc.org/sqlite"
)

type config struct {
	port          int
	dsn           string
	runMigrations bool
}

type application struct {
	config         config
	models         *models.Models
	sessionManager *scs.SessionManager
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "HTTP Server port")
	flag.StringVar(&cfg.dsn, "dsn", "database.db", "Database Service Name")
	flag.BoolVar(&cfg.runMigrations, "run-migrations", false, "Run database migrations")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := dbConnect(cfg.dsn)
	if err != nil {
		slog.Error("Unable to stablish database connection", "err", err)
		return
	}

	if cfg.runMigrations {
		slog.Info("Running database migrations")
		if err := models.RunMigrations(db); err != nil {
			slog.Error("Unable to run database migrations", "err", err)
			return
		}
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true

	app := &application{
		models:         models.New(db),
		config:         cfg,
		sessionManager: sessionManager,
	}

	s := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", cfg.port),
		Handler: app.routes(),
	}

	slog.Info("Starting web server", "port", cfg.port)
	err = s.ListenAndServe()
	slog.Error("Something went wrong with the server", "error", err)
}

func dbConnect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
