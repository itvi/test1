package main

import (
	"ams/pkg/models/mysqlite"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	_ "github.com/mattn/go-sqlite3"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	errorLog        *log.Logger
	infoLog         *log.Logger
	session         *sessions.Session
	templateCache   map[string]*template.Template
	users           *mysqlite.UserModel
	assets          *mysqlite.AssetModel
	assetCategory   *mysqlite.AssetCategoryModel
	assetStatus     *mysqlite.AssetStatusModel
	assetManagement *mysqlite.AssetManagementMode
	computer        *mysqlite.ComputerModel
}

func main() {
	addr := flag.String("addr", "localhost:8000", "HTTP network address")
	flag.Parse()

	// custom logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret")

	// database
	db, err := openDB("./db.db")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	app := &application{
		errorLog:        errorLog,
		infoLog:         infoLog,
		session:         session,
		templateCache:   templateCache,
		users:           &mysqlite.UserModel{DB: db},
		assets:          &mysqlite.AssetModel{DB: db},
		assetCategory:   &mysqlite.AssetCategoryModel{DB: db},
		assetStatus:     &mysqlite.AssetStatusModel{DB: db},
		assetManagement: &mysqlite.AssetManagementMode{DB: db},
		computer:        &mysqlite.ComputerModel{DB: db},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 12 * time.Second, // 原来是10秒，导致“net::ERR_EMPTY_RESPONSE”。默认是多少？
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// database
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
