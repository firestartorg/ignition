package monitor

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/firestart/ignition/x/application"
	"net/http"
	"time"
)

var (
	HookHealth    application.Hook = "health"
	HookReadiness application.Hook = "health/ready"
	HookLiveness  application.Hook = "health/live"

	ErrNotReady = errors.New("application is not ready")
)

// WithMonitor adds a health monitor to the application.
func WithMonitor(opts ...Option) application.Option {
	// Apply the options
	options := &Options{
		// Default
		port:    8080,
		timeout: time.Second * 5,
	}
	for _, opt := range opts {
		opt(options)
	}

	return func(app application.App, hooks *application.Hooks) {
		//var listener *net.Listener
		router := httprouter.New()

		// Add the metrics monitor
		if options.metrics != nil {
			router.Handler(http.MethodGet, "/metrics", options.metrics)
		}
		// Add the readiness monitor
		if options.readiness {
			router.GET("/health/ready", readinessProbeHandle(app, hooks, options))
		}
		// Add the liveness monitor
		if options.liveness {
			router.GET("/health/live", livenessProbeHandle(app, hooks, options))
		}

		// Add a startup hook
		hooks.AddStartup(func(ctx context.Context, app application.App) error {
			err := http.ListenAndServe(fmt.Sprintf(":%d", options.port), router)
			return err
		})

		//// Add a shutdown hook
		//hooks.AddShutdown(func(ctx context.Context, app application.App) error {
		//	return (*listener).Close()
		//})
	}
}

func readinessProbeHandle(app application.App, hooks *application.Hooks, options *Options) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ok := getStatus(context.Background(), app, hooks, HookHealth, options) &&
			getStatus(context.Background(), app, hooks, HookReadiness, options)

		if !ok {
			http.Error(w, "Unavailable", http.StatusServiceUnavailable)
		} else {
			http.Error(w, "Ok", http.StatusOK)
		}
	}
}

func livenessProbeHandle(app application.App, hooks *application.Hooks, options *Options) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ok := getStatus(context.Background(), app, hooks, HookHealth, options) &&
			getStatus(context.Background(), app, hooks, HookLiveness, options)

		if !ok {
			http.Error(w, "Unavailable", http.StatusServiceUnavailable)
		} else {
			http.Error(w, "Ok", http.StatusOK)
		}
	}
}

// getStatus returns the status of the given health hook
func getStatus(ctx context.Context, app application.App, hooks *application.Hooks, hook application.Hook, options *Options) bool {
	pctx, cancel := context.WithTimeout(ctx, options.timeout)
	defer cancel()

	// Check all the hooks
	err := hooks.RunWithContext(hook, app, pctx)
	if err != nil {
		return false
	}
	return true
}
