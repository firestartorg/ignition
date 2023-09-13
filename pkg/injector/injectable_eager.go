package injector

type EagerInjectable[T any] struct {
	instance T
}

func (e EagerInjectable[T]) get(injector *Injector) (T, error) {
	return e.instance, nil
}

func newEagerInjectable[T any](instance T) *EagerInjectable[T] {
	return &EagerInjectable[T]{
		instance: instance,
	}
}

var _ Injectable[any] = (*EagerInjectable[any])(nil)
