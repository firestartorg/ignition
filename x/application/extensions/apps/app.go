package apps

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/monitor"
	"runtime"
)

var (
	// AppName is the name of the injectable that contains the app globals
	AppName = "app"

	// AppVersion is the gauge metric exposed
	AppVersion = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "go",
			Subsystem: "app",
			Name:      "info",
			Help:      "Metric with a constant '1' value labeled by version and goversion from which the cloud gateway receiver was built.",
		},
		[]string{"app_name", "app_version", "go_version"},
	)
)

type appGlobals struct {
	name    string
	version string
}

// WithVersion sets the application version
func WithVersion(name, version string) application.Option {
	return func(app application.App, hooks *application.Hooks) {
		// Inject app globals
		injector.InjectNamed(app.Injector, AppName, appGlobals{name, version})

		// Add the following metric
		AppVersion.WithLabelValues(name, version, runtime.Version()).Set(1)
	}
}

type HealthCheckFunc = application.HookFunc

// WithHealthCheck adds a health check to the application.
func WithHealthCheck(f HealthCheckFunc) application.Option {
	return func(app application.App, hooks *application.Hooks) {
		hooks.Add(monitor.HookHealth, f)
	}
}

// WithAppReadinessCheck adds a readiness check to the application.
func WithReadinessCheck(f HealthCheckFunc) application.Option {
	return func(app application.App, hooks *application.Hooks) {
		hooks.Add(monitor.HookReadiness, f)
	}
}

// WithAppLivenessCheck adds a liveness check to the application.
func WithLivenessCheck(f HealthCheckFunc) application.Option {
	return func(app application.App, hooks *application.Hooks) {
		hooks.Add(monitor.HookLiveness, f)
	}
}
