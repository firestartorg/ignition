package hook

type Collection interface {
	// Hooks returns all hooks
	Hooks() []any
	// AddHook adds a hook to the Hooks.
	AddHook(hook any)
}

type collection struct {
	// hooks is a slice of all hooks that have been registered.
	hooks []any
}

func (h *collection) Hooks() []any {
	return h.hooks
}

func (h *collection) AddHook(hook any) {
	// Add the hook to the slice of hooks
	h.hooks = append(h.hooks, hook)
}

// NewCollection creates a new Collection.
func NewCollection() Collection {
	return &collection{}
}
