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
