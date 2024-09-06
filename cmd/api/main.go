package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/tiwanakd/Calculator-API/internal/models"
)

type api struct {
	logger        *slog.Logger
	calculations  *models.CalculationModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":8000", "port to run the web server on")
	dsn := flag.String("dsn", "postgres://api:pass@localhost/calculatordb?sslmode=disable", "PostgreSQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	api := &api{
		logger:        logger,
		calculations:  &models.CalculationModel{DB: db},
		templateCache: templateCache,
	}

	api.logger.Info("server running at port", "addr", *addr)

	err = http.ListenAndServe(*addr, api.routes())
	api.logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
