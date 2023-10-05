package presets

import (
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/http"
	"gitlab.com/firestart/ignition/x/application/extensions/sentry/tracing"
)

// WithHttpServer adds the http server to the application
func WithHttpServer() application.Option {
	return http.WithServer(
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
				WithHttpServer(),
			},
		)...,
	)
}
