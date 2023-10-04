package http

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/x/application"
	"net/http"
)

func newMiddleware(handler http.Handler, app application.App, hooks *application.Hooks) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Server", "Ignition")

		// Get the request context
		ctx := r.Context()

		// Recover from panics
		defer func() {
			if err := recover(); err != nil {
				log.Ctx(ctx).Panic().Interface("panic", err).Msg("Recovered from panic")

				// Process the request context
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("{\"error\":\"Internal Server Error\"}"))
			}
		}()

		// If there are no hooks, just serve the request
		if hooks == nil {
			handler.ServeHTTP(w, r)
			return
		}

		var err error
		// Process the request context
		ctx, err = hooks.ProcessContext(application.HookRequest, ctx, app)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to process request context")
			return
		}

		// Process the request with the new context
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
