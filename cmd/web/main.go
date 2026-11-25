package main

import (
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"letsgo.rumi/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lmittmann/tint"
)

type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	devMode        bool
}

func main() {
	// Configs.
	addr := getEnv("WEB_ADDR", ":4000")
	databaseUrl := getEnv("WEB_DATABASE_URL", "root:very-hard-password@/mydb?parseTime=true")
	devMode := getEnv("WEB_DEV_MODE", "true") == "true"
	cookieSecure := getEnv("WEB_COOKIE_SECURE", "false") == "true"

	// Our logger.
	logger := slog.New(tint.NewTextHandler(os.Stdout, nil))

	// Log if dev mode is active.
	if devMode {
		logger.Info("dev mode is active")
	}

	// Connect to database.
	db, err := openDB(databaseUrl)
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

	// Form decoder.
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = cookieSecure

	// Initialize a new instance of application struct, containing all the
	// dependencies.
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		devMode:        devMode,
	}

	// Setup http server.
	srv := &http.Server{
		Addr:     addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", addr)
	err = srv.ListenAndServe()
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
