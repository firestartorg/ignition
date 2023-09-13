package injector

type Provider[T any] func(*Injector) (T, error)
