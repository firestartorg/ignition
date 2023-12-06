package main

import "gitlab.com/firestart/ignition/pkg/injector"

func main() {
	// This is a simple example of how to inject a dependencies.

	// First we need to create a new injector.
	// If you use the application package, you can use the following code:
	// inj := application.Injector
	inj := injector.NewInjector()

	// Then we can either register a new static dependency or a dynamic dependency.
	// A static dependency is a dependency that is always the same.
	// A dynamic dependency is a dependency that is created the first time it is requested.

	// So let's register a static dependency first.
	injector.Inject(inj, staticDependency{})

	// Now we can register a dynamic dependency.
	// Be aware, that unnamed dependencies do not support injection of interface types. See: examples/injector/injector_named
	injector.Provide(inj, provideDynamicDependency)

	// If one would like to get either a static or dynamic dependency, one can use the following code:
	// dep, err := injector.Get[type](inj)
	// For example:
	dep, err := injector.Get[dynamicDependency](inj)

	_ = dep
	_ = err
}

type staticDependency struct {
	// This is a static dependency
}

type dynamicDependency struct {
	// This is a dynamic dependency
	dep staticDependency
}

func provideDynamicDependency(inj *injector.Injector) (dynamicDependency, error) {
	// This is the function that will be called when the dynamic dependency is requested.
	// The injector will pass itself as the first argument.

	// First we need to get the static dependency from the injector.
	dep, err := injector.Get[staticDependency](inj)
	if err != nil {
		return dynamicDependency{}, err
	}

	return dynamicDependency{
		dep: dep,
	}, nil
}
