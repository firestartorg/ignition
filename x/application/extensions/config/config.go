package config

import (
	"gitlab.com/firestart/ignition/pkg/config"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
)

// WithSettings sets the settings
func WithSettings() application.Option {
	return func(app application.App, hooks *application.Hooks) {
		dyn, err := config.LoadConfig()
		if err != nil {
			panic(err)
		}

		injector.Inject(app.Injector, *dyn)
	}
}
