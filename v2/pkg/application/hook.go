package application

type StartupHook interface {
	Startup() error
}

type ShutdownHook interface {
	Shutdown() error
}

type startupHook struct {
	f func() error
}

func (h startupHook) Startup() error {
	return h.f()
}

type shutdownHook struct {
	f func() error
}

func (h shutdownHook) Shutdown() error {
	return h.f()
}

// NewStartupHook creates a new StartupHook.
func NewStartupHook(f func() error) StartupHook {
	return startupHook{f}
}

// NewShutdownHook creates a new ShutdownHook.
func NewShutdownHook(f func() error) ShutdownHook {
	return shutdownHook{f}
}
