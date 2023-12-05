package config

import (
	"gitlab.com/firestart/ignition/pkg/config"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
)

// WithConfiguration sets the settings
func WithConfiguration() application.Option {
	return func(app application.App) {
		dyn, err := config.LoadConfig()
		if err != nil {
			panic(err)
		}

		injector.Inject(app.Injector, *dyn)
	}
}
