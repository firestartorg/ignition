package http

import (
	"github.com/julienschmidt/httprouter"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
	"net/http"
)

// AddNamedRoute adds a route to the named http server
func AddNamedRoute(app application.App, serverName string, method string, path string, handle httprouter.Handle) error {
	// Get the server container
	srv, err := injector.GetNamed[server](app.Injector, serverName)
	if err != nil {
		return err
	}

	// Add the route
	srv.router.Handle(method, path, handle)

	return nil
}

// MustAddNamedRoute adds a route to the named http server
// and panics if an error occurs
func MustAddNamedRoute(app application.App, serverName string, method string, path string, handle httprouter.Handle) {
	err := AddNamedRoute(app, serverName, method, path, handle)
	if err != nil {
		panic(err)
	}
}

// AddRoute adds a route to the default http server
func AddRoute(app application.App, method string, path string, handle httprouter.Handle) error {
	return AddNamedRoute(app, ServerName, method, path, handle)
}

// MustAddRoute adds a route to the default http server
// and panics if an error occurs
func MustAddRoute(app application.App, method string, path string, handle httprouter.Handle) {
	MustAddNamedRoute(app, ServerName, method, path, handle)
}

// MustAddGetRoute adds a GET route to the default http server
// and panics if an error occurs
func MustAddGetRoute(app application.App, path string, handle httprouter.Handle) {
	MustAddNamedRoute(app, ServerName, http.MethodGet, path, handle)
}

// MustAddPostRoute adds a POST route to the default http server
// and panics if an error occurs
func MustAddPostRoute(app application.App, path string, handle httprouter.Handle) {
	MustAddNamedRoute(app, ServerName, http.MethodPost, path, handle)
}

// MustAddPutRoute adds a PUT route to the default http server
// and panics if an error occurs
func MustAddPutRoute(app application.App, path string, handle httprouter.Handle) {
	MustAddNamedRoute(app, ServerName, http.MethodPut, path, handle)
}

// MustAddPatchRoute adds a PATCH route to the default http server
// and panics if an error occurs
func MustAddPatchRoute(app application.App, path string, handle httprouter.Handle) {
	MustAddNamedRoute(app, ServerName, http.MethodPatch, path, handle)
}

// MustAddDeleteRoute adds a DELETE route to the default http server
// and panics if an error occurs
func MustAddDeleteRoute(app application.App, path string, handle httprouter.Handle) {
	MustAddNamedRoute(app, ServerName, http.MethodDelete, path, handle)
}
