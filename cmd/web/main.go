package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/tsamba120/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

// define an application struct to inject log dependencies into our handlers
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// network address
	addr := flag.String("addr", ":4000", "HTTP network address")
	// db conn string
	// must be either localhost or name of container or DNS of database service
	dsn := flag.String("dsn", "web:1234@tcp(go-mysql:3306)/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	// set up logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// instantiate a new http.Server struct but attached out custom error logger
	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(), // call app.routes method
		ErrorLog: errorLog,
	}

	// Use http.ListenAndServe() to state a new webserver. Pass in TCP network address to listen on and servemux
	infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) { // multi-return value funcs have r.v.'s on right side
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
