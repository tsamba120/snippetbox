package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	// network address
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// set up logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Use the http.NewServeMux() function to initialize a new ServeMux (aka router)
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// create file server that serves files out of "./ui/static/ directory"
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// instantiate a new http.Server struct but attached out custom error logger
	srv := &http.Server{
		Addr:     *addr,
		Handler:  mux,
		ErrorLog: errorLog,
	}

	// Use http.ListenAndServe() to state a new webserver. Pass in TCP network address to listen on and servemux
	infoLog.Printf("Starting server on %s\n", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
