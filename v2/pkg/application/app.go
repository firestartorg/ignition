package application

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/v2/pkg/hook"
	"gitlab.com/firestart/ignition/v2/pkg/inject"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	// Injector is the injector used by the application
	inject.Injector
	// Hooks is the hooks used by the application
	hook.Collection
}

// New creates a new App
func New(opts ...Extension) App {
	// Create the application
	app := App{
		Injector:   inject.NewInjector(),
		Collection: hook.NewCollection(),
	}

	// Apply the options
	for _, opt := range opts {
		opt(&app)
	}

	return app
}

// Shutdown shuts down the application
func (a *App) Shutdown() {
	log.Info().Msg("Shutting down application")

	hooks := hook.Gather[ShutdownHook](a)
	for _, h := range hooks {
		err := h.Shutdown()
		if err != nil {
			log.Error().Err(err).Msg("Failed to run shutdown hooks")
		}
	}

	log.Info().Msg("Stopped application")
}

// Run runs the application
func (a *App) Run() {
	log.Info().Msg("Starting application")

	// Shutdown the application when the run function returns
	defer a.Shutdown()

	// Listen for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Run the startup hooks
	done := make(chan struct{})
	go func() {
		// Gather all startup hooks
		hooks := hook.Gather[StartupHook](a)

		var wg sync.WaitGroup
		wg.Add(len(hooks))

		// Run all startup hooks
		for _, h := range hooks {
			go func(h StartupHook) {
				defer wg.Done()
				err := h.Startup()
				if err != nil {
					log.Error().Err(err).Msg("Failed to run startup hooks")
				}
			}(h)
		}

		wg.Wait()

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
