package main

import (
	"gitlab.com/firestart/ignition/v2/extension/config"
	"gitlab.com/firestart/ignition/v2/pkg/application"
	"gitlab.com/firestart/ignition/v2/pkg/inject"
)

func main() {
	app := application.New(
		config.WithConfig(
			config.WithName("name-of-config"), // set the name of the injectable
			config.WithPreload(),
		),
	)
	app.AddHook(&lifecycleManager{})

	inject.ProvideValue(app, lifecycleManager{})
}

type lifecycleManager struct {
}

func (l lifecycleManager) Shutdown() error {
	//TODO implement me
	panic("implement me")
}

func (l lifecycleManager) Startup() error {
	//TODO implement me
	panic("implement me")
}

var _ application.StartupHook = (*lifecycleManager)(nil)
var _ application.ShutdownHook = (*lifecycleManager)(nil)
