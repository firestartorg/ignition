package injector

type Injectable[T any] interface {
	get(*Injector) (T, error)

	//healthCheck(ctx context.Context) error
	//shutdown() error
}
