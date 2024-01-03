# Application

Application is a package that provides a set of functions to manage the application lifecycle. 
Its designed to be flexible and extensible.

## Usage

Usage of the application package is simple. Create a new application and run it.

```go
package main

import (
  "gitlab.com/firestart/ignition/x/application"
)

func main() {
  app := application.New()
  app.Run()
}
```

Except for the most basic applications, you will want to add some extensions to the application.
For example, if you want to use the configuration extension, you can add it like this:
  
```go
package main

import (
  "gitlab.com/firestart/ignition/x/application"
  "gitlab.com/firestart/ignition/x/application/extensions/config"
)

func main() {
  app := application.New(
    // Add the configuration extension
    config.WithConfiguration(),
  )
  app.Run()
}
```

For more information on the extensions available, see the [extensions](#extensions) section.

## Extensions

Extensions are packages that provide additional functionality to the application. Ignition 
provides a number of extensions that can be used to add functionality to your application.

- [grpc](extensions/grpc/README.md) - Provides grpc server and client functionality
- [http](extensions/http/README.md) - Provides http server functionality
- [monitor](extensions/monitor/README.md) - Provides monitoring functionality using [prometheus](https://github.com/prometheus/client_golang)
- config - Provides configuration management
- logger - Provides logging functionality using [zerolog](https://github.com/rs/zerolog)
- sentry - Provides error reporting using [sentry](https://sentry.io)
- sentry/recovery - Provides recovery using [sentry](https://sentry.io)
- sentry/tracing - Provides tracing using [sentry](https://sentry.io)
- vcs - Provides version control system information
- apps - Makes applications FireStart conformant

Ignition also provides a number of presets that can be used to quickly add functionality to your application. 
For more information on presets, see the [presets](#presets) section.

### How to write an extension

Extensions are simple to write. They are just functions that take an application. Therefore they have access to
all of the application's functionality, including the injector and hooks.

```go
package myextension

import (
  "gitlab.com/firestart/ignition/x/application"
)

func WithMyExtension() application.Option {
  return func(app application.Application) {
    // Do something here
  }
}
```

For example, the http extension is written like this:

```go
package http

import (
  "gitlab.com/firestart/ignition/pkg/injector"
  "gitlab.com/firestart/ignition/x/application"
)

func WithHttpServer() application.Option {
  return func(app application.Application) {
    server := newHttpServer()
    // Add the server to the injector
    injector.InjectNamed(app.Injector, "http-server", server)
    // Start the server when the application starts
    app.AddHook(application.HookStartup, func(ctx context.Context, app App) error {
      return server.Start()
    })
    // Stop the server when the application stops
    app.AddHook(application.HookShutdown, func(ctx context.Context, app App) error {
      return server.Stop()
    })
  }
}
```

## Presets

Presets are a collection of extensions that can be used to quickly add functionality to your application.

- [blank](#blank-preset) - Base application with vcs, config, logger, sentry, monitor, and apps extensions
- [grpc](#grpc-preset) - Blank application with grpc extension (optional client)
- [http](#http-preset) - Blank application with http extension

### Blank Preset

The blank preset is the base application with a number of extensions. It is the default preset.

```go
package main

import (
  "gitlab.com/firestart/ignition/x/application/presets"
)

func main() {
  app := presets.NewBlankApp(
    "example-blank-app", // The name of the application
    // Add additional extensions here, if needed
  )
  app.Run()
}
```

### Grpc Preset

The grpc preset is the blank preset with the grpc extension. It also includes the grpc client preset.

```go
package main

import (
  "gitlab.com/firestart/ignition/x/application/presets"
)

func main() {
  app := presets.NewRpcApp(
    "example-grpc-app", // The name of the application
    // Add additional extensions here, if needed
    // If you want to use the grpc client preset, add it here, like this:
    presets.WithRpcClientFactory(),
  )
  app.Run()
}
```

For more information on grpc servers and clients, see the [grpc](extensions/grpc/README.md) extension.

### Http Preset

The http preset is the blank preset with the http extension.

```go
package main

import (
  "gitlab.com/firestart/ignition/x/application/presets"
)

func main() {
  app := presets.NewHttpApp(
    "example-http-app", // The name of the application
    // Add additional extensions here, if needed
  )
  app.Run()
}
```

For more information on http servers, see the [http](extensions/http/README.md) extension.
