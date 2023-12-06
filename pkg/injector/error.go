package injector

type ErrorType int

const (
	ErrTypeUnknown ErrorType = iota
	ErrTypeInjectableNotFound
)

type Error struct {
	// Type of error
	Type ErrorType
	// Description of error like for instance the name of the injectable
	Descriptor string
}

func (e *Error) Error() string {
	switch e.Type {
	case ErrTypeInjectableNotFound:
		return "injectable not found: " + e.Descriptor

	case ErrTypeUnknown:
		fallthrough
	default:
		return "unknown error"
	}
}
