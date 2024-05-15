package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {

	// mux is of type serveMux which is a struct in the http package which has methods to
	// process route requests sent to the app's http server
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Use the mux.HandleFunc() method to register the routes and their handlers
	// Note the handlers are called as methods on the 'app' configuration struct
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	return mux
}
