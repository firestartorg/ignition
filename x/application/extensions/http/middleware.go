package http

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/x/application"
	"net/http"
)

// Middleware is a function that wraps a http handler
type Middleware = func(handler http.Handler) http.Handler

// newBaseMiddleware creates a new the default middleware for the server
func newBaseMiddleware(handler http.Handler, app application.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Server", "Ignition")

		// Get the request context
		ctx := r.Context()

		// If there are no hooks, just serve the request
		if app.Hooks == nil {
			handler.ServeHTTP(w, r)
			return
		}

		var err error
		// Process the request context
		ctx, err = app.Hooks.ProcessContext(application.HookRequest, ctx, app)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to process request context")
			return
		}

		// Process the request with the new context
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
