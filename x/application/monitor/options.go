package monitor

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/firestart/ignition/x/application"
	"net/http"
	"time"
)

type Options struct {
	port int

	metrics http.Handler

	readiness bool
	liveness  bool
	timeout   time.Duration
}

type Option = func(opts *Options)

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

// WithTimeout sets the timeout for the monitor.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.timeout = timeout
	}
}
