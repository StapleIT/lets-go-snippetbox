package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// the struct "application" is instantiated in func Main() in the main.go file of the same package (main)
// as the variable 'app' which contains config data (e.g. a pointer to the logger function)
// the handlers are created as methods on this struct type
// The reason for this is so that values in the struct (e.g. the logger function)
// can be accessed from within each handler. In func main(), the handler is called as app.handlername
// this means the logger comes along for the ride and is available to the handler method function

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
	//w.Write([]byte("Hellow from Snippetbox"))
}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// msg := fmt.Sprintf("This is a view of snippet with ID: %d", id)

	// w.Write([]byte(msg))

	fmt.Fprintf(w, "This is a view of snippet with ID: %d", id)
}
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is snippet create"))
}
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
