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

	// HookInit is a hook that is called when the application is initialized.
	// This hook is primarily used for context processing.
	HookInit Hook = "init"
)

type HookFunc func(ctx context.Context, app App) error
type ContextHookFunc func(ctx context.Context, app App) (context.Context, error)

// Hooks is a struct that can be used to inject hooks into the application
type Hooks struct {
	hooks map[Hook][]HookFunc
	// hooksMutex is a mutex to protect the hooks map from concurrent access
	hooksMutex sync.Mutex

	// contextHooks is a map of context hooks
	contextHooks map[Hook][]ContextHookFunc
	// contextHooksMutex is a mutex to protect the contextHooks map from concurrent access
	contextHooksMutex sync.Mutex
}

func newHooks() *Hooks {
	return &Hooks{
		hooks:             make(map[Hook][]HookFunc),
		hooksMutex:        sync.Mutex{},
		contextHooks:      map[Hook][]ContextHookFunc{},
		contextHooksMutex: sync.Mutex{},
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

// AddContext adds a context hook to the application
func (h *Hooks) AddContext(hook Hook, f ContextHookFunc) {
	h.contextHooksMutex.Lock()
	defer h.contextHooksMutex.Unlock()

	// Create the hook map if it doesn't exist
	if _, ok := h.contextHooks[hook]; !ok {
		h.contextHooks[hook] = make([]ContextHookFunc, 0)
	}

	h.contextHooks[hook] = append(h.contextHooks[hook], f)
}

// ProcessContext runs the context hooks for the given hook
func (h *Hooks) ProcessContext(hook Hook, ctx context.Context, app App) (context.Context, error) {
	h.contextHooksMutex.Lock()
	defer h.contextHooksMutex.Unlock()

	// Create the hook map if it doesn't exist
	if _, ok := h.contextHooks[hook]; !ok {
		return ctx, nil
	}

	for _, f := range h.contextHooks[hook] {
		var err error
		ctx, err = f(ctx, app)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

// RunWithContext runs the hooks for the given hook
func (h *Hooks) RunWithContext(hook Hook, app App, ctx context.Context) error {
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
	ctx, err := h.context(app)
	if err != nil {
		return err
	}

	return h.RunWithContext(hook, app, ctx)
}

func (h *Hooks) shutdown(app App) error {
	h.hooksMutex.Lock()
	defer h.hooksMutex.Unlock()

	// Create the hook map if it doesn't exist
	if _, ok := h.hooks[HookShutdown]; !ok {
		return nil
	}

	// Create a context
	ctx, err := h.context(app)
	if err != nil {
		return err
	}

	// loop through the hooks in reverse order
	for i := len(h.hooks[HookShutdown]) - 1; i >= 0; i-- {
		f := h.hooks[HookShutdown][i]
		err = f(ctx, app)
		if err != nil {
			return err
		}
	}

	return nil
}

// waitUntil runs the hooks and waits until they are all done
func (h *Hooks) waitUntil(hook Hook, app App) error {
	ctx, err := h.context(app)
	if err != nil {
		return err
	}

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

// context returns the context for the given hook
func (h *Hooks) context(app App) (context.Context, error) {
	return h.ProcessContext(HookInit, context.Background(), app)
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
