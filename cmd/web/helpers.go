package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/nosurf"
	"github.com/yunushamod/snippetbox/pkg/models"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name.
	// If no entry exists in the cache with the provided name, call the serverError
	// helper method that we made earlier.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)
	// Write the template to the buffer, instead of straight to the http.ResponseWriter
	// If there's an error, call our serverError helper and return.
	// Execute the template set, passing in any dynamic data.
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.CurrentYear = time.Now().Year()
	td.AuthenticatedUser = app.authenticatedUser(r)
	// Add the flash message to the template data if one exists.
	td.Flash = app.session.PopString(r, "flash")
	return td
}

func (app *application) authenticatedUser(r *http.Request) *models.User {
	//return app.session.GetInt(r, "userID")
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}
