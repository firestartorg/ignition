package presets

import (
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
)

func pack(opts ...[]application.Option) []application.Option {
	var locals []application.Option
	for _, opt := range opts {
		locals = append(locals, opt...)
	}
	return locals
}

type Configure[T interface{}] func(config T) []application.Option

// MakeConfigurable is a helper function to make a list of options configurable
func MakeConfigurable[T interface{}](key string, configure Configure[T]) application.Option {
	return func(app application.App, hooks *application.Hooks) {
		config, _ := injector.ExtractConfig[T](app.Injector, key)

		options := configure(config)
		for _, option := range options {
			option(app, hooks)
		}
	}
}
