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

// RunWithContext runs the hooks for the given hook
func (h *Hooks) RunWithContext(hook Hook, ctx context.Context, app App) error {
	h.hooksMutex.Lock()
	defer h.hooksMutex.Unlock()

	// Create the hook map if it doesn't exist
	if _, ok := h.hooks[hook]; !ok {
		return nil
	}

	for _, f := range h.hooks[hook] {
		if err := f(ctx, app); err != nil {
			return err
		}
	}

	return nil
}

// Run runs the hooks for the given hook
func (h *Hooks) Run(hook Hook, app App) error {
	return h.RunWithContext(hook, context.Background(), app)
}

// waitUntil runs the hooks and waits until they are all done
func (h *Hooks) waitUntil(hook Hook, app App) error {
	ctx := context.Background()

	fs := h.cloneHook(hook)

	// Create a wait group
	var wg sync.WaitGroup
	// Add the number of hooks to the wait group
	wg.Add(len(fs))

	for _, f := range fs {
		go func(f HookFunc) {
			defer wg.Done()

			err := f(ctx, app)
			if err != nil {
				app.logger.Error().Err(err).Msg("Hook failed")
				panic(err)
			}
		}(f)
	}

	wg.Wait()

	return nil
}

func (h *Hooks) cloneHook(hook Hook) []HookFunc {
	h.hooksMutex.Lock()
	defer h.hooksMutex.Unlock()

	// Create the hook map if it doesn't exist
	if _, ok := h.hooks[hook]; !ok {
		return []HookFunc{}
	}

	// Copies the hooks
	hooks := make([]HookFunc, len(h.hooks[hook]))
	copy(hooks, h.hooks[hook])

	return hooks
}
