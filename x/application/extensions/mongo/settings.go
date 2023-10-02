package mongo

import "go.mongodb.org/mongo-driver/mongo/options"

type settings struct {
	// Config is a flag that indicates that the injected Config should be used
	Config bool

	// ClientOptions is the options to use when connecting to the database
	ClientOptions *options.ClientOptions
	// Database is the default database to use
	Database string
	// Collections is a map of collections names to their actual names
	Collections map[string]string
}

// newSettings creates a new settings struct with default values
func newSettings(opts ...Option) settings {
	opt := settings{
		Database:    "firestart",
		Collections: make(map[string]string),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

type Option func(*settings)

// WithDatabase sets the Database name
func WithDatabase(database string) Option {
	return func(o *settings) {
		o.Database = database
	}
}

// WithClientOptions sets the client settings
func WithClientOptions(clientOptions *options.ClientOptions) Option {
	return func(o *settings) {
		o.ClientOptions = clientOptions
	}
}

// WithNamedCollection sets a Collections collection
func WithNamedCollection(name, collection string) Option {
	return func(o *settings) {
		o.Collections[name] = collection
	}
}

// FromSettings sets the Config flag
func FromSettings() Option {
	return func(o *settings) {
		o.Config = true
	}
}
