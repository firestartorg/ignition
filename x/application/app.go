package application

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/pkg/injector"
)

type Option = func(app App, hooks *Hooks)

type App struct {
	// Injector is the injector used by the application
	*injector.Injector

	// logger is the logger used by the application
	logger zerolog.Logger
	// hooks is the hooks used by the application
	hooks *Hooks
}

// New creates a new App
func New(opts ...Option) App {
	// Create the application
	app := App{
		Injector: injector.NewInjector(),
		hooks:    newHooks(),

		logger: log.With().Timestamp().Logger(),
	}

	// Apply the options
	for _, opt := range opts {
		opt(app, app.hooks)
	}

	return app
}

// Run runs the application
func (a App) Run() {
	a.logger.Info().Msg("Starting application")
	defer func() {
		a.logger.Info().Msg("Stopping application")

		// Run the shutdown hooks
		err := a.hooks.Run(HookShutdown, a)
		if err != nil {
			panic(err)
		}
	}()

	// Run the startup hooks
	err := a.hooks.waitUntil(HookStartup, a)
	if err != nil {
		panic(err)
	}
}
