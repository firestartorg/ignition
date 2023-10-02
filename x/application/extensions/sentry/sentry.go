package sentry

import (
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
)

// WithSentry adds a sentry client to the application.
func WithSentry(opts ...Option) application.Option {
	s := newSettings(opts...)

	return func(app application.App, hooks *application.Hooks) {
		if s.config {
			var err error
			s, err = injector.ExtractConfig[settings](app.Injector, "Sentry")
			if err != nil {
				panic(err)
			}
			// Override the config with the specified options
			for _, opt := range opts {
				opt(&s)
			}
		}
	}
}

// WithDefaultSentry adds a sentry client to the application.
func WithDefaultSentry() application.Option {
	return WithSentry(FromSettings())
}
