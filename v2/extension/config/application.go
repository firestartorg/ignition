package config

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/v2/pkg/application"
	"gitlab.com/firestart/ignition/v2/pkg/config"
	"gitlab.com/firestart/ignition/v2/pkg/inject"
)

const (
	// DefaultConfigName is the default name of the configuration file.
	DefaultConfigName = "ignition/config"
)

type Option func(*options)

type options struct {
	// Preload indicates whether the configuration should be preloaded before the application starts.
	preload bool
	// configName is the configName of the configuration configName to load.
	configName string
	// isAbsolute indicates whether the configuration configName is an absolute path.
	isAbsolute bool
	// name is the name of the injected configuration.
	name string
}

func WithConfig(opts ...Option) application.Extension {
	// Set default options
	o := options{}
	// Apply options
	for _, opt := range opts {
		opt(&o)
	}

	return func(app *application.App) {
		// If the configuration should be preloaded, load it now
		if o.preload {
			c, err := loadConfig(o)
			if err != nil {
				log.Fatal().Err(err).Str("name", o.name).Msg("Failed to preload configuration")
				return
			}
			inject.ProvideNamedValue(app, getInjectableName(o), c)
			return
		}

		// Otherwise, load the configuration when the application starts
		app.AddHook(application.NewStartupHook(func() error {
			c, err := loadConfig(o)
			if err != nil {
				log.Fatal().Err(err).Str("name", o.name).Msg("Failed to load configuration")
				return err
			}
			inject.ProvideNamedValue(app, getInjectableName(o), c)
			return nil
		}))
	}
}

func loadConfig(o options) (*config.Config, error) {
	if o.configName != "" {
		if o.isAbsolute {
			return config.LoadFile(o.configName)
		}
		return config.Load(o.configName)
	}
	return config.LoadConfig()
}

func getInjectableName(o options) string {
	if o.name != "" {
		return o.name
	}
	return DefaultConfigName
}

// WithPreload indicates whether the configuration should be preloaded before the application starts.
func WithPreload() Option {
	return func(opts *options) {
		opts.preload = true
	}
}

// WithConfigFile sets the configName of the configuration configName to load.
func WithConfigFile(file string) Option {
	return func(opts *options) {
		opts.configName = file
		opts.isAbsolute = true
	}
}

// WithConfigName sets the configName of the configuration configName to load.
func WithConfigName(name string) Option {
	return func(opts *options) {
		opts.configName = name
	}
}

// WithName sets the name of the injected configuration.
func WithName(name string) Option {
	return func(opts *options) {
		opts.name = name
	}
}
