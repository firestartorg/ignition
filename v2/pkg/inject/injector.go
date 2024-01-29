package inject

import (
	"errors"
	"sync"
)

var (
	ErrUnknownInjectable = errors.New("unknown injectable")
)

type Injector interface {
	// SetInjectable an injectable by type name or name
	SetInjectable(name string, isTypeName bool, injectable any)
	// GetInjectable an injectable by type name or name
	GetInjectable(name string, isTypeName bool) (any, error)
}

type injector struct {
	// mu is a mutex to protect the Injector from concurrent access
	mu sync.RWMutex

	// namedInjectables are injectables that are named
	namedInjectables map[string]interface{}
	// injectables are injectables that are not named
	injectables map[string]interface{}
}

func (i *injector) SetInjectable(name string, isTypeName bool, injectable any) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if isTypeName {
		i.injectables[name] = injectable
	} else {
		i.namedInjectables[name] = injectable
	}
}

func (i *injector) GetInjectable(name string, isTypeName bool) (any, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var injectable any
	var ok bool
	if isTypeName {
		injectable, ok = i.injectables[name]
	} else {
		injectable, ok = i.namedInjectables[name]
	}
	if !ok {
		return nil, &Error{Type: ErrTypeInjectableNotFound, Descriptor: name}
	}

	return injectable, nil
}

func NewInjector() Injector {
	return &injector{
		mu:               sync.RWMutex{},
		namedInjectables: make(map[string]interface{}),
		injectables:      make(map[string]interface{}),
	}
}
