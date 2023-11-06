package logging

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
)

func WithZerolog(logger zerolog.Logger) application.Option {
	return func(app application.App) {
		hookFn := func(ctx context.Context, app application.App) (context.Context, error) {
			return logger.WithContext(ctx), nil
		}

		// Set the global logger
		log.Logger = logger
		zerolog.DefaultContextLogger = &logger // Also set the default context logger for zerolog (used if log.Ctx(ctx) is called)

		app.AddContextProcessor(application.HookInit, hookFn)
		app.AddContextProcessor(application.HookRequest, hookFn)
	}
}

// WithConfigurableZerolog adds a zerolog logger to the application that can be configured via the config.
func WithConfigurableZerolog(logger zerolog.Logger) application.Option {
	return func(app application.App) {
		conf, err := injector.ExtractConfig[loggingConfig](app.Injector, "App:Logging")
		if err != nil {
			return
		}
		if conf.Level == "" {
			WithZerolog(logger.Level(zerolog.WarnLevel))(app)
			return
		}

		level, err := zerolog.ParseLevel(conf.Level)
		if err != nil {
			WithZerolog(logger.Level(zerolog.WarnLevel))(app)
			return
		}

		WithZerolog(logger.Level(level))(app)
	}
}

type loggingConfig struct {
	Level string
}
