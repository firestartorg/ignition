package injector

import "gitlab.com/firestart/ignition/pkg/config"

// InjectConfig injects the configuration into the injector
func InjectConfig[T any](injector *Injector) error {
	dyn, err := config.LoadConfig()
	if err != nil {
		return err
	}

	var conf T
	err = dyn.Unpack(&conf)
	if err != nil {
		return err
	}

	Inject(injector, *dyn)
	Inject(injector, conf)
	return nil
}

// MustInjectConfig injects the configuration into the injector, panicking if an error occurs
func MustInjectConfig[T any](injector *Injector) {
	err := InjectConfig[T](injector)
	if err != nil {
		panic(err)
	}
}

// ExtractConfig extracts a sub configuration from the injected configuration
func ExtractConfig[T any](injector *Injector, key string) (value T, err error) {
	var conf config.Config
	conf, err = Get[config.Config](injector)
	if err != nil {
		return
	}

	err = conf.Sub(key).Unpack(&value)
	if err != nil {
		return
	}

	return
}
