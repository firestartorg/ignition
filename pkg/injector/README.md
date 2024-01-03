# Injector

The injector package provides a simple dependency injection system, that makes use of go generics to provide a simple and type safe way to inject dependencies into your application.

## Usage

First you need to create an injector. This is done by calling the `NewInjector` function (or use the injector provided by
the application instance e.g. `app.Injector`). *Note: The injector is thread safe.* 

```go
inj := injector.NewInjector()
```

Next you need to add some dependencies to the injector. This is done by calling the `Provide` function. The `Provide` function takes a function (so called provider) that:

- receives the injector as a parameter (so it can access other dependencies)
- returns the dependency and an error (if the dependency could not be created)

```go
injector.Provide(inj, func(inj *injector.Injector) (Service, error) {
  return Service{}, nil
})
```

Or if you want to provide a dependency that is already created:

```go
injector.Inject(inj, Service{})
```

Both of these methods can be used to provide typed dependencies. However, if you want to provide an **interface type**, or 
the same type multiple times, you need to use the `ProvideNamed` function or the `InjectNamed` function. These functions
take a additional string parameter that is used to identify the dependency.

```go
injector.ProvideNamed(inj, "service", func(inj *injector.Injector) (Service, error) {
  return Service{}, nil
})

injector.InjectNamed(inj, "service", Service{})
```

Once you have provided all of your dependencies, you can use the `Get`, `GetNamed`, `MustGet` and `MustGetNamed` functions
to retrieve your dependencies.

```go
service, err := injector.Get[Service](inj, Service{})
service, err := injector.GetNamed[Service](inj, "service", Service{})
service := injector.MustGet[Service](inj, Service{})
service := injector.MustGetNamed[Service](inj, "service", Service{})
```
