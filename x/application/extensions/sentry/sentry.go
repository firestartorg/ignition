package sentry

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
	"time"
)

// WithSentry adds a sentry client to the application.
func WithSentry(opts ...Option) application.Option {
	s := newSettings(opts...)

	return func(app application.App) {
		// If the config flag is set, extract the config from the injector
		if s.config {
			var err error
			s, err = injector.ExtractConfig[Settings](app.Injector, "Sentry")
			if err != nil {
				panic(err)
			}
			// Override the config with the specified options
			for _, opt := range opts {
				opt(&s)
			}
		}
		if !s.Enabled {
			log.Warn().Msg("Sentry is disabled")
			return
		}

		// Initialize sentry
		err := sentry.Init(sentry.ClientOptions{
			Dsn:         s.Dsn,
			Environment: s.Environment,
			Release:     s.Release,

			EnableTracing:    s.EnableTracing,
			TracesSampleRate: s.TracesSampleRate,
			Debug:            s.Debug,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize sentry")
		}

		// AddHook a sentry hub to the application context for each request
		app.AddContextProcessor(application.HookRequest, func(ctx context.Context, app application.App) (context.Context, error) {
			hub := sentry.GetHubFromContext(ctx)
			if hub == nil {
				hub = sentry.CurrentHub().Clone()
				ctx = sentry.SetHubOnContext(ctx, hub)
			}

			return ctx, nil
		})

		// AddHook a shutdown hook to flush sentry events
		app.AddShutdownHook(func(ctx context.Context, app application.App) error {
			// Flush buffered events before the program terminates.
			// Set the timeout to the maximum duration the program can afford to wait.
			sentry.Flush(time.Duration(s.FlushTimeout) * time.Second)

			return nil
		})
	}
}

// WithDefaultSentry adds a sentry client to the application.
func WithDefaultSentry() application.Option {
	return WithSentry(FromSettings())
}
