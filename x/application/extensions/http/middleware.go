package http

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/x/application"
	"net/http"
)

func newMiddleware(handler http.Handler, app application.App, hooks *application.Hooks) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Server", "Ignition")

		// If there are no hooks, just serve the request
		if hooks == nil {
			handler.ServeHTTP(w, r)
			return
		}

		var err error
		// Process the request context
		ctx := r.Context()
		ctx, err = hooks.ProcessContext(application.HookRequest, ctx, app)
		if err != nil {
			log.Error().Err(err).Msg("Failed to process request context")
			return
		}

		// Process the request with the new context
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
