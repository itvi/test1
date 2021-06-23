package main

import (
	"ams/pkg/models/mysqlite"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	assets        *mysqlite.AssetModel
	assetCategory *mysqlite.AssetCategoryModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// custom logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	// database
	db, err := openDB("./db.db")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		assets:        &mysqlite.AssetModel{DB: db},
		assetCategory: &mysqlite.AssetCategoryModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
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
