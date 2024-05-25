package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// convert the return type to http.Handler from *http.ServeMux
func (app *application) routes() http.Handler {

	// this print function was only called once when app was started
	// fmt.Println("Calling the routes function")

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

	// Pass the servemux as the 'next' parameter to the commonHeaders middleware.
	// Because commonHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.
	//return app.recoverPanic(app.logRequest(commonHeaders(mux)))

	// refactor to use 'alice' to manage middleware chaining
	// create a new 'standard' alice.Chain containing our previous 'before' request
	// middleware functions
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)
}
