package logging

import (
	"context"
	"github.com/rs/zerolog"
	"gitlab.com/firestart/ignition/x/application"
	"os"
)

func WithZerolog() application.Option {
	return func(app application.App, hooks *application.Hooks) {
		hooks.AddContext(application.HookInit, func(ctx context.Context, app application.App) (context.Context, error) {
			return zerolog.New(os.Stderr).With().Timestamp().Stack().Caller().Logger().WithContext(ctx), nil
		})
	}
}
