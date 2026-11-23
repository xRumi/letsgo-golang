package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"letsgo.rumi/internal/models"

	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lmittmann/tint"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
	devMode       bool
}

func main() {
	// Flags, make environment variables default if available.
	addr := flag.String("addr", getEnv("WEB_ADDR", ":4000"), "HTTP network address")
	dsn := flag.String("dsn", getEnv("WEB_DSN", "root:very-hard-password@/mydb?parseTime=true"), "MySQL data source name")
	devMode := flag.Bool("devMode", false, "Developer mode")
	flag.Parse()

	// Our logger.
	logger := slog.New(tint.NewTextHandler(os.Stdout, nil))

	// Log if dev mode is active.
	if *devMode {
		logger.Info("dev mode is active")
	}

	// Connect to database.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connected")

	// Cache templates.
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	// Initialize a new instance of application struct, containing all the
	// dependencies.
	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
		devMode:       *devMode,
	}

	logger.Info("starting server", "port", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// openDB() opens MySQL database connection.
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
