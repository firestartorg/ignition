package application

import (
	"context"
	"sync"
)

type Hook = string

const (
	// HookStartup is a hook that is called when the application starts
	HookStartup Hook = "startup"
	// HookShutdown is a hook that is called when the application shuts down
	HookShutdown Hook = "shutdown"
)

type HookFunc func(ctx context.Context, app App) error

// Hooks is a struct that can be used to inject hooks into the application
type Hooks struct {
	hooks map[Hook][]HookFunc
	// hooksMutex is a mutex to protect the hooks map from concurrent access
	hooksMutex sync.Mutex
}

func newHooks() *Hooks {
	return &Hooks{
		hooks:      make(map[Hook][]HookFunc),
		hooksMutex: sync.Mutex{},
	}
}

// Add adds a hook to the application
func (h *Hooks) Add(hook Hook, f HookFunc) {
	h.hooksMutex.Lock()
	defer h.hooksMutex.Unlock()

	// Create the hook map if it doesn't exist
	if _, ok := h.hooks[hook]; !ok {
		h.hooks[hook] = make([]HookFunc, 0)
	}

	h.hooks[hook] = append(h.hooks[hook], f)
}

// AddStartup adds a startup hook to the application
func (h *Hooks) AddStartup(f HookFunc) {
	h.Add(HookStartup, f)
}

// AddShutdown adds a shutdown hook to the application
func (h *Hooks) AddShutdown(f HookFunc) {
	h.Add(HookShutdown, f)
}

// runWithContext runs the hooks for the given hook
func (h *Hooks) runWithContext(hook Hook, ctx context.Context, app App) error {
	h.hooksMutex.Lock()
	defer h.hooksMutex.Unlock()

	// Create the hook map if it doesn't exist
	if _, ok := h.hooks[hook]; !ok {
		return nil
	}

	for _, f := range h.hooks[hook] {
		go func(f HookFunc) {
			if err := f(ctx, app); err != nil {
				panic(err)
			}
		}(f)
	}

	return nil
}

// run runs the hooks for the given hook
func (h *Hooks) run(hook Hook, app App) error {
	return h.runWithContext(hook, context.Background(), app)
}
