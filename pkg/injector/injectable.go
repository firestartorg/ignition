package injector

type Injectable[T any] interface {
	get(*Injector) (T, error)
}
