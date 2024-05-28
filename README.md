# lets-go-snippetbox

This is the code created while following along to the book "Let's Go" by Alex Edwards.

There was no attempt to tag each version of each file when it was refactored or modified during the course.  The version here is the latest version at any point of time depending on where I was up to in the book.

## Major conceptual refactors:

### Project Structure (section 2.7):  
 - splitting handler code out of main.go into handlers.go
 - creating directory structure with cmd, internal and ui forming the first layer
 - handlers.go and main.go are in a web folder under cmd folder

### HTML Templating (section 2.8):
 - concept of using html/template package to render ftml templates from handlers

### Convert Handlers to Methods on an Application Struct (section 3.3)
 - create an app struct in main.go which contains dependencies which need to be accessed by multiple handlers (e.g. our structured logger)
 - because we are developing the app in one package, handlers in a different file can still be coded as methods on the struct defined in main.go
 - func main() in main.go (the handlers' calling function) is refactored as is every handler in handlers.go

### Centralised errors and the helpers.go (section 3.4)
 - helpers.go (used throughout the app) added to the web folder
 - centralised error handling functions (e.g. serverError) created as methods on the app struct
 - app.serverError() called in handler functions

### Move routes out of main.go into routes.go (section 3.5)
 - new routes.go file created in web folder
 - app.routes() method created which handles the mapping of routes to handlers via mux.HandleFunc as before
 - note that this method returns a *http.Servemux type which slots in perfectly to replace "mux" in the http.ListenAndServe() function in func main()

### Database and snippets.go (chapter 4)
 - create snippets.go in internal folder as part of a new 'models package
 - snippets.go contains methods on the db model which pexecute SQL statements to create or read snippets to from the database
 - include the snippet model in the app struct in main.go
 - modify handlers to use the database model to create or read snippet information

### Dynamic Templates (section 5.1)
 - itroduced html/template functionality to parse template files and run functions in {{.variable}} form

### Template Cache and Template Errors (section 5.2 and 5.3)
 - cache helper function to map parsed templates and add them to the app dependency struct to be used by the handlers when needed
 - simplify the handlers to use a new render() function which renders templates from the cached map of parsed templates.
 - the render() function is adapted to first render to a buffer, check if an error occurs and only render to client if no error triggered, otherwise raise our serverError.
 - add a custom dynamic template function (humanDate)

### Middleware! (chapter 6)
 - used like 'before request' or 'after request' in some frameworks
 - for 'before all requests' the midleware must happen before the mux. To do this we must change the route.go signature to http.Handler from *http.ServeMux and return the middleware instead of mux.
 - refactor to use 'alice' library which chains middleware functions more easily.  It wasn't that difficult before since it was just a case of each middleware calling the next before, but 'alice' can build a slice of handlers and automatically chain them.

 ### Forms (chapter 7)
 - passing form data with validation
 - refactor validation to use helpers for case where there are many forms and a lot of validation
   - new concept!  the validator.Validator struct is 'embedded' in the snippetCreateForm struct which therefore inherits all the fields and methods of the validator.Validator struct!!
 - refactor to use 3rd party library "go-playground/form" as a 'form decoder' to automatically parse forms data into our struct!  So you build the struct to represent the data (as we had already) and the 'form' library decodes the form data into it.  The intent is that, with large or many forms, it reduces the amount of code to be written to perform the r.PostForm.Get() on each form field.

 ### Sessions (chapter 8)
 - add functionality of 3rd party session store to share data between requests

