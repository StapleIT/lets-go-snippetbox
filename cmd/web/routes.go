package main

import (
	"net/http"

	"github.com/StapleIT/lets-go-snippetbox/ui"

	"github.com/justinas/alice"
)

// convert the return type to http.Handler from *http.ServeMux
func (app *application) routes() http.Handler {

	// this print function was only called once when app was started
	// fmt.Println("Calling the routes function")

	// mux is of type serveMux which is a struct in the http package which has methods to
	// process route requests sent to the app's http server
	mux := http.NewServeMux()

	// Use the http.FileServerFS() function to create a HTTP handler which
	// serves the embedded files in ui.Files. It's important to note that our
	// static files are contained in the "static" folder of the ui.Files
	// embedded filesystem. So, for example, our CSS stylesheet is located at
	// "static/css/main.css". This means that we no longer need to strip the
	// prefix from the request URL -- any requests that start with /static/ can
	// just be passed directly to the file server and the corresponding static
	// file will be served (so long as it exists).
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.Handle("GET /ping", http.HandlerFunc(ping))
	// Use the mux.HandleFunc() method to register the routes and their handlers
	// Note the handlers are called as methods on the 'app' configuration struct

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

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

	// add handlers for user registration and login
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	// refactor to use 'alice' to manage middleware chaining
	// create a new 'standard' alice.Chain containing our previous 'before' request
	// middleware functions
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)
}
