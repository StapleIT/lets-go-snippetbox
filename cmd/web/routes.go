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

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Update these routes to use the new dynamic middleware chain followed by
	// the appropriate handler function. Alice's Then() and ThenFunc() are the only
	// two 'then' functions the libraru has.  Then() takes in a http.Handler type and
	// returns a http.Handler; ThenFunc() takes in a http.HandlerFunc and ALSO returns
	// a http.Handler type. But it seems that it also coerces the handler function
	// to http.Handler as it is happy to take it in without error!
	// Note that because the alice ThenFunc()
	// method returns a http.Handler, we now no longer register it
	// with the mux.HandleFunc, we register it with the mux.Handle().

	mux.Handle("GET /{$}", dynamic.Then(http.HandlerFunc(app.home))) // note that this strategy works just as well as ThenFunc below
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// refactor to use 'alice' to manage middleware chaining
	// create a new 'standard' alice.Chain containing our previous 'before' request
	// middleware functions
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)
}
