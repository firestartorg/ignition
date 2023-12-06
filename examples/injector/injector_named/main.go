package main

import "gitlab.com/firestart/ignition/pkg/injector"

func main() {
	// This is a simple example of how to inject a named dependencies
	// or an interface dependency

	// First we need to create a new injector.
	// If you use the application package, you can use the following code:
	// inj := application.Injector
	inj := injector.NewInjector()

	// Then we can either register a new static dependency or a dynamic dependency.
	// A static dependency is a dependency that is always the same.
	// A dynamic dependency is a dependency that is created the first time it is requested.

	// So let's register a static dependency first.
	injector.InjectNamed(inj, "greeting", "Hello World!")

	// Now we can register a dynamic dependency that implements the greeter interface.
	injector.ProvideNamed(inj, DepName, provideGreeter)

	// Now we can register a dynamic dependency.
	injector.ProvideNamed(inj, "dynamic", provideDynamicDependency)
}

const DepName = "greeter"

type greeter interface {
	Greet() string
}

type greeterImpl struct {
	greeting string
}

func (g greeterImpl) Greet() string {
	return g.greeting
}

func provideGreeter(inj *injector.Injector) (greeter, error) {
	greeting, err := injector.GetNamed[string](inj, "greeting")
	if err != nil {
		return nil, err
	}

	return &greeterImpl{
		greeting: greeting,
	}, nil
}

type dynamicDependency struct {
}

func provideDynamicDependency(_ *injector.Injector) (dynamicDependency, error) {
	return dynamicDependency{}, nil
}
