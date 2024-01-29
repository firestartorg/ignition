package config

import (
	"gitlab.com/firestart/ignition/v2/pkg/config"
	"gitlab.com/firestart/ignition/v2/pkg/inject"
)

// GetConfig returns the unpacked configuration from the injector with the given key.
// If the key is empty, the root configuration is used.
func GetConfig[T any](injector inject.Injector, key string) (value T, err error) {
	var conf *config.Config
	conf, err = inject.Get[*config.Config](injector)
	if err != nil {
		return
	}
	// If the key is empty, use the root configuration
	if key != "" {
		conf = conf.Sub(key)
	}
	// Unpack the configuration
	err = conf.Unpack(&value)
	if err != nil {
		return
	}
	return
}
