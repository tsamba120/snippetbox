package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// define an application struct to inject log dependencies into our handlers
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// network address
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:1234@tcp(localhost:3306)/snippetbox?parseTime=true", "MySQL data source name")
	// dsn := flag.String("dsn", "root:password@tcp(127.0.0.1:3306)/test", "MySQL data source name")

	flag.Parse()

	// set up logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

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
