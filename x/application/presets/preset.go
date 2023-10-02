package presets

import "gitlab.com/firestart/ignition/x/application"

func pack(opts ...[]application.Option) []application.Option {
	var locals []application.Option
	for _, opt := range opts {
		locals = append(locals, opt...)
	}
	return locals
}
