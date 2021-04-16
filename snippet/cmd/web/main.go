package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"snippet/pkg/models/postgres"

	_ "github.com/lib/pq"
)

var (
	dbcon = "host=localhost port=5432 user=postgres password=password dbname=snippets"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *postgres.SnippetModel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDB(dbcon)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	// Initialize a mysql.SnippetModel instance and add it to the application
	// dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &postgres.SnippetModel{DB: db},
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
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
