package application

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/pkg/injector"
	"os"
	"os/signal"
	"syscall"
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

// Shutdown shuts down the application
func (a App) Shutdown() {
	log.Info().Msg("Shutting down application")

	// Run the shutdown hooks
	err := a.hooks.shutdown(a)
	if err != nil {
		log.Error().Err(err).Msg("Failed to run shutdown hooks")
		return
	}

	log.Info().Msg("Stopped application")
}

// Run runs the application
func (a App) Run() {
	log.Info().Msg("Starting application")

	// Shutdown the application when the run function returns
	defer a.Shutdown()

	// Listen for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Run the startup hooks
	done := make(chan struct{})
	go func() {
		// Run the startup hooks
		err := a.hooks.waitUntil(HookStartup, a)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to run startup hooks")
		}

		done <- struct{}{}
	}()

	// Wait for shutdown or app to be done
	select {
	case sig := <-shutdown:
		log.Info().Msgf("Received %s signal", sig)
		return
	case <-done:
		return
	}
}
