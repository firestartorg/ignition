package monitor

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
	"net"
	"net/http"
)

var (
	HookHealth    application.Hook = "health"
	HookReadiness application.Hook = "health/ready"
	HookLiveness  application.Hook = "health/live"

	ErrNotReady = errors.New("application is not ready")
)

// WithMonitor adds a health monitor to the application.
func WithMonitor(opts ...Option) application.Option {
	return func(app application.App, hooks *application.Hooks) {
		// Apply the options
		options := newOptions(opts...)
		// Check if the config should be used
		if options.config {
			conf, err := injector.ExtractConfig[Config](app.Injector, "App:Monitor")
			if err == nil {
				conf.apply(&options)
			}
		}

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
			// Start the server
			err := http.ListenAndServe(fmt.Sprintf(":%d", options.port), router)
			// Check if the error is a "port already in use" error
			var opErr *net.OpError
			if options.portAuto && errors.As(err, &opErr) && opErr.Op == "listen" && opErr.Net == "tcp" {
				// Try again
				var listener net.Listener
				listener, err = net.Listen("tcp", ":0")
				if err != nil {
					return err
				}
				port := listener.Addr().(*net.TCPAddr).Port

				// Log the port change
				log.Ctx(ctx).Info().Int("port", port).Msg("monitor: port already in use, using available port")

				// Start the server
				return http.Serve(listener, router)
			} else {
				return err
			}
		})

		//// Add a shutdown hook
		//hooks.AddShutdown(func(ctx context.Context, app application.App) error {
		//	return (*listener).Close()
		//})
	}
}

func readinessProbeHandle(app application.App, hooks *application.Hooks, options Options) httprouter.Handle {
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

func livenessProbeHandle(app application.App, hooks *application.Hooks, options Options) httprouter.Handle {
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
func getStatus(ctx context.Context, app application.App, hooks *application.Hooks, hook application.Hook, options Options) bool {
	pctx, cancel := context.WithTimeout(ctx, options.timeout)
	defer cancel()

	// Check all the hooks
	err := hooks.RunWithContext(hook, app, pctx)
	if err != nil {
		return false
	}
	return true
}
