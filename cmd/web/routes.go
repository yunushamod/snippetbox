package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

// http.Handler instead of *http.ServeMux
func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	//standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//mux := http.NewServeMux()
	mux := pat.New()
	mux.Get("/", app.session.Enable(noSurf(app.authenticate(http.HandlerFunc((app.home))))))
	mux.Get("/snippet/create", app.session.Enable(noSurf(app.authenticate(app.requireAuthenticatedUser(http.HandlerFunc(app.createSnippetForm))))))
	mux.Post("/snippet/create", app.session.Enable(noSurf(app.authenticate(app.requireAuthenticatedUser(http.HandlerFunc(app.createSnippet))))))
	mux.Get("/snippet/:id", app.session.Enable(noSurf(app.authenticate(http.HandlerFunc(app.showSnippet)))))

	// Add five new routes.
	mux.Get("/user/signup", app.session.Enable(noSurf(app.authenticate(app.authenticate(http.HandlerFunc(app.signupUserForm))))))
	mux.Post("/user/signup", app.session.Enable(noSurf(app.authenticate(http.HandlerFunc(app.signupUser)))))
	mux.Get("/user/login", app.session.Enable(noSurf(app.authenticate(http.HandlerFunc(app.loginUserForm)))))
	mux.Post("/user/login", app.session.Enable(noSurf(app.authenticate(http.HandlerFunc(app.loginUser)))))
	mux.Post("/user/logout", app.session.Enable(noSurf(app.authenticate(app.requireAuthenticatedUser(http.HandlerFunc(app.logoutUser))))))
	mux.Get("/ping", http.HandlerFunc(ping))
	/*
		mux.HandleFunc("/", app.home)
		mux.HandleFunc("/snippet", app.showSnippet)
		mux.HandleFunc("/snippet/create", app.createSnippet)
	*/
	// Create a file server which serves files out of the "./ui/static" directories
	// Note that the path given to the http.Dir function is relative to the project directory root
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler
	// all URL paths that start with "/static/". For matching paths, we strip the "/static"
	// prefix before the request reaches the file server.
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	//return standardMiddleware.Then(mux)
}
