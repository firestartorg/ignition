package sentry

type Settings struct {
	// config is a flag that indicates that the injected config should be used
	config bool

	// FlushTimeout is the timeout for the sentry flush. Defaults to 2s.
	FlushTimeout int

	// Dsn is the sentry DSN
	Dsn              string
	EnableTracing    bool
	TracesSampleRate float64
	Debug            bool
}

// newSettings creates a new Settings struct
func newSettings(opts ...Option) Settings {
	s := Settings{
		FlushTimeout:     2,
		TracesSampleRate: 1,
	}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

type Option func(*Settings)

// FromSettings sets the config flag
func FromSettings() Option {
	return func(s *Settings) {
		s.config = true
	}
}

// WithDsn sets the sentry DSN
func WithDsn(dsn string) Option {
	return func(s *Settings) {
		s.Dsn = dsn
	}
}

// WithDebugFlag sets the sentry debug flag
func WithDebugFlag(debug bool) Option {
	return func(s *Settings) {
		s.Debug = debug
	}
}

// WithTracesSampleRate sets the sentry traces sample rate
func WithTracesSampleRate(tracesSampleRate float64) Option {
	return func(s *Settings) {
		s.TracesSampleRate = tracesSampleRate
	}
}

// WithFlushTimeout sets the sentry flush timeout
func WithFlushTimeout(flushTimeout int) Option {
	return func(s *Settings) {
		s.FlushTimeout = flushTimeout
	}
}

// EnableTracing sets the sentry enable tracing flag
func EnableTracing() Option {
	return func(s *Settings) {
		s.EnableTracing = true
	}
}
