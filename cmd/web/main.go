package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"letsgo.rumi/internal/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lmittmann/tint"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	// The value of flag will be stored in the addr variable at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "root:very-hard-password@/mydb?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(tint.NewTextHandler(os.Stdout, nil))
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connected")

	// Initialize a new instance of application struct, containing
	// the dependencies
	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("starting server", "port", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
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
