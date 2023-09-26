package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.com/firestart/ignition/x/application"
	"runtime"
)

var (
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

// WithAppVersion sets the application version
func WithAppVersion(name, version string) application.Option {
	return func(app application.App, hooks *application.Hooks) {
		AppVersion.WithLabelValues(name, version, runtime.Version()).Set(1)
	}
}
