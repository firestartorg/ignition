package monitor

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/firestart/ignition/x/application"
	"net/http"
	"time"
)

type Options struct {
	port     int
	portAuto bool // If true, the port will be automatically assigned if it is already in use

	metrics http.Handler

	readiness bool
	liveness  bool
	timeout   time.Duration

	// If true, the injected config should be used
	config bool
}

type Option = func(opts *Options)

func newOptions(opts ...Option) Options {
	o := Options{
		port:    8080,
		timeout: 5 * time.Second,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithDefaultMonitor adds a default health monitor to the application.
func WithDefaultMonitor() application.Option {
	return WithMonitor(
		WithPrometheusMetrics(),
		WithReadiness(),
		WithLiveness(),
	)
}

// WithMetricsHandler adds a metrics monitor to the application.
func WithMetricsHandler(handler http.Handler) Option {
	return func(opts *Options) {
		opts.metrics = handler
	}
}

// WithPrometheusMetrics adds a prometheus metrics monitor to the application.
func WithPrometheusMetrics() Option {
	return WithMetricsHandler(promhttp.Handler())
}

// WithReadiness adds a readiness monitor to the application.
func WithReadiness() Option {
	return func(opts *Options) {
		opts.readiness = true
	}
}

// WithLiveness adds a liveness monitor to the application.
func WithLiveness() Option {
	return func(opts *Options) {
		opts.liveness = true
	}
}

// WithPort sets the port for the monitor.
func WithPort(port int) Option {
	return func(opts *Options) {
		opts.port = port
	}
}

// WithPortAntiCollision sets the port for the monitor to be automatically assigned if it is already in use.
func WithPortAntiCollision() Option {
	return func(opts *Options) {
		opts.portAuto = true
	}
}

// WithTimeout sets the timeout for the monitor.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.timeout = timeout
	}
}

// FromConfig sets the config flag.
func FromConfig() Option {
	return func(opts *Options) {
		opts.config = true
	}
}
