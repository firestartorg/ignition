package inject

import "sync"

type LazyInjectable[T any] struct {
	mu sync.Mutex

	// instance is the instance of the injectable
	instance T

	// provider is the provider of the injectable
	provider Provider[T]
	// provided is true if the instance has been provided
	provided bool
}

// nolint:unused
func (l *LazyInjectable[T]) get(injector *Injector) (T, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.provided {
		instance, err := l.provider(injector)
		if err != nil {
			var empty T
			return empty, err
		}

		l.instance = instance
		l.provided = true
	}

	return l.instance, nil
}

func newLazyInjectable[T any](provider Provider[T]) *LazyInjectable[T] {
	return &LazyInjectable[T]{
		mu:       sync.Mutex{},
		provider: provider,
		provided: false,
	}
}

var _ Injectable[any] = (*LazyInjectable[any])(nil)

type Provider[T any] func(*Injector) (T, error)
