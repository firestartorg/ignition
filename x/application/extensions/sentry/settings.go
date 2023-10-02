package sentry

type settings struct {
	// config is a flag that indicates that the injected config should be used
	config bool
}

// newSettings creates a new settings struct
func newSettings(opts ...Option) settings {
	s := settings{}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

type Option func(*settings)

// FromSettings sets the config flag
func FromSettings() Option {
	return func(s *settings) {
		s.config = true
	}
}
