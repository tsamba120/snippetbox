package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tsamba120/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

// define an application struct to inject log dependencies into our handlers
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// network address
	addr := flag.String("addr", ":4000", "HTTP network address")
	// db conn string
	// must be either localhost or name of container or DNS of database service
	// dsn := flag.String("dsn", "web:1234@tcp(go-mysql:3306)/snippetbox?parseTime=true", "MySQL data source name")
	dsn := flag.String("dsn", "web:1234@tcp(localhost:3306)/snippetbox?parseTime=true", "MySQL data source name")
	// session secret (random key used to encrypt and authenticate session cookies -- 32 bytes long)
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")

	flag.Parse()

	// set up logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// connect to db
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	// initialize new session manager, passing in secret key
	// configured to expire every 12 hours
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	// create application struct that we inject into server
	app := application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
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
