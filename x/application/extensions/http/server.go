// Http router extension for applications

package http

import (
	"context"
	"errors"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
	"net/http"
)

var (
	// ServerName is the name of the injectable that contains the router
	ServerName = "ignition/http-router"
)

// server is a container for the server
type server struct {
	// options is the server options
	options ServerOptions
	// router is the router
	router *httprouter.Router

	// server is the listener
	server *http.Server
}

func WithServer(opts ...ServerOption) application.Option {
	return WithNamedServer(ServerName)
}

// WithNamedServer adds a named http router to the application.
func WithNamedServer(name string, opts ...ServerOption) application.Option {
	// Create the server options with the defaults
	options := newServerOptions(opts...)

	return func(app application.App, hooks *application.Hooks) {
		// Inject the server container into the application
		injector.InjectNamed(app.Injector, name, server{
			options: options,
			router:  httprouter.New(),
		})

		// Add a startup hook
		hooks.AddStartup(func(ctx context.Context, app application.App) error {
			// Get the server container
			srv, err := injector.GetNamed[server](app.Injector, name)
			if err != nil {
				return err
			}

			// Wrap the router with the middlewares
			handler := srv.options.wrapMiddleware(srv.router)
			handler = newBaseMiddleware(handler, app, hooks)

			// Create the server
			srv.server = &http.Server{Addr: srv.options.addr(), Handler: handler}
			injector.InjectNamed(app.Injector, name, srv) // Update the server container

			err = srv.server.ListenAndServe()
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		})

		// Add a shutdown hook
		hooks.AddShutdown(func(ctx context.Context, app application.App) error {
			// Get the server container
			srv, err := injector.GetNamed[server](app.Injector, name)
			if err != nil {
				return err
			}

			return srv.server.Shutdown(ctx)
		})
	}
}
