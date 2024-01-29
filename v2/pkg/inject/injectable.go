package inject

type Injectable[T any] interface {
	get(Injector) (T, error)
}
