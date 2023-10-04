package logging

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/x/application"
	"os"
)

func WithZerolog(logger zerolog.Logger) application.Option {
	return func(app application.App, hooks *application.Hooks) {
		hookFn := func(ctx context.Context, app application.App) (context.Context, error) {
			return logger.WithContext(ctx), nil
		}

		// Set the global logger
		log.Logger = logger

		hooks.AddContext(application.HookInit, hookFn)
		hooks.AddContext(application.HookRequest, hookFn)
	}
}

func WithDefaultZerolog() application.Option {
	return WithZerolog(zerolog.New(os.Stderr).With().Timestamp().Stack().Caller().Logger())
}
