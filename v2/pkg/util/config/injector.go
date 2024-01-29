package configutil

import (
	"gitlab.com/firestart/ignition/v2/pkg/config"
	"gitlab.com/firestart/ignition/v2/pkg/inject"
)

// InjectConfig injects the configuration into the injector
func InjectConfig[T any](injector inject.Injector) error {
	dyn, err := config.LoadConfig()
	if err != nil {
		return err
	}

	var conf T
	err = dyn.Unpack(&conf)
	if err != nil {
		return err
	}

	inject.ProvideValue(injector, *dyn)
	inject.ProvideValue(injector, conf)
	return nil
}

// MustInjectConfig injects the configuration into the injector, panicking if an error occurs
func MustInjectConfig[T any](injector inject.Injector) {
	err := InjectConfig[T](injector)
	if err != nil {
		panic(err)
	}
}

// ExtractConfig extracts a sub configuration from the injected configuration
func ExtractConfig[T any](injector inject.Injector, key string) (value T, err error) {
	var conf config.Config
	conf, err = inject.Get[config.Config](injector)
	if err != nil {
		return
	}

	err = conf.Sub(key).Unpack(&value)
	if err != nil {
		return
	}

	return
}
