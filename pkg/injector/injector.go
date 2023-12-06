package injector

import (
	"errors"
	"sync"
)

var (
	ErrUnknownInjectable = errors.New("unknown injectable")
)

type Injector struct {
	// mu is a mutex to protect the Injector from concurrent access
	mu sync.RWMutex

	// namedInjectables are injectables that are named
	namedInjectables map[string]interface{}
	// injectables are injectables that are not named
	injectables map[string]interface{}
}

func NewInjector() *Injector {
	return &Injector{
		mu:               sync.RWMutex{},
		namedInjectables: make(map[string]interface{}),
		injectables:      make(map[string]interface{}),
	}
}

func (i *Injector) set(name string, injectable any, named bool) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if named {
		i.namedInjectables[name] = injectable
	} else {
		i.injectables[name] = injectable
	}
}

func (i *Injector) get(name string, named bool) (any, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var injectable any
	var ok bool
	if named {
		injectable, ok = i.namedInjectables[name]
	} else {
		injectable, ok = i.injectables[name]
	}
	if !ok {
		return nil, &Error{Type: ErrTypeInjectableNotFound, Descriptor: name}
	}

	return injectable, nil
}
