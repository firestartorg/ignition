package injector

import (
	"reflect"
)

var defaultInjector = NewInjector()

func getOrDefaultInjector(injector *Injector) *Injector {
	if injector == nil {
		injector = defaultInjector
	}
	return injector
}

func getTypeName[T any](value T) string {
	valType := reflect.TypeOf(value)
	if valType.Kind() == reflect.Ptr {
		valType = valType.Elem()
	}
	return valType.PkgPath() + "/" + valType.Name()
}

// Inject injects the given value into the injector
func Inject[T any](injector *Injector, value T) {
	injector = getOrDefaultInjector(injector)
	injector.set(getTypeName(value), newEagerInjectable(value), false)
}

// InjectNamed injects the given value into the injector with the given name
func InjectNamed[T any](injector *Injector, name string, value T) {
	injector = getOrDefaultInjector(injector)
	injector.set(name, newEagerInjectable(value), true)
}

// Provide injects the given provider into the injector
func Provide[T any](injector *Injector, provider Provider[T]) {
	injector = getOrDefaultInjector(injector)
	var value T
	injector.set(getTypeName(value), newLazyInjectable(provider), false)
}

// ProvideNamed injects the given provider into the injector with the given name
func ProvideNamed[T any](injector *Injector, name string, provider Provider[T]) {
	injector = getOrDefaultInjector(injector)
	injector.set(name, newLazyInjectable(provider), true)
}

func get[T any](injector *Injector, name string, named bool) (T, error) {
	injector = getOrDefaultInjector(injector)

	var value T
	injableAny, err := injector.get(name, named)
	if err != nil {
		return value, err
	}

	injable, ok := injableAny.(Injectable[T])
	if !ok {
		return value, ErrUnknownInjectable
	}

	return injable.get(injector)
}

// Get returns the value of the given type from the injector
func Get[T any](injector *Injector) (T, error) {
	var value T
	return get[T](injector, getTypeName(value), false)
}

// GetNamed returns the value of the given type from the injector with the given name
func GetNamed[T any](injector *Injector, name string) (T, error) {
	return get[T](injector, name, true)
}

// MustGet returns the value of the given type from the injector, panicking if an error occurs
func MustGet[T any](injector *Injector) T {
	value, err := Get[T](injector)
	if err != nil {
		panic(err)
	}
	return value
}

// MustGetNamed returns the value of the given type from the injector with the given name, panicking if an error occurs
func MustGetNamed[T any](injector *Injector, name string) T {
	value, err := GetNamed[T](injector, name)
	if err != nil {
		panic(err)
	}
	return value
}
