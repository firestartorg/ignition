package presets

import (
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/http"
	"gitlab.com/firestart/ignition/x/application/extensions/sentry/tracing"
)

// WithHttpServer adds the http server to the application
func WithHttpServer(port int16) application.Option {
	return http.WithServer(
		http.WithPort(port),
		http.WithMiddleware(tracing.HttpMiddleware),
	)
}

// NewHttpApp creates a new application with the following components:
// - blank app preset (see NewBlankApp)
func NewHttpApp(name string, opts ...application.Option) application.App {
	return NewBlankApp(
		name,
		pack(
			opts,
			[]application.Option{
				MakeConfigurable("App", func(config httpConfig) []application.Option {
					// SetInjectable the default port
					if config.Port == 0 {
						config.Port = 3000
					}

					return []application.Option{
						WithHttpServer(config.Port),
					}
				}),
			},
		)...,
	)
}

type httpConfig struct {
	// Port is the port to listen on
	Port int16
}
