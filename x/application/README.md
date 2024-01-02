# Application

Application is a package that provides a set of functions to manage the application lifecycle. 
Its designed to be flexible and extensible.

## Usage

Usage of the application package is simple. Create a new application and run it.

```go
package main

import (
  "gitlab.com/firestart/ignition/pkg/application"
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
  "gitlab.com/firestart/ignition/pkg/application"
  "gitlab.com/firestart/ignition/pkg/application/extensions/config"
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

## Presets

Presets are a collection of extensions that can be used to quickly add functionality to your application.

- blank - Base application with vcs, config, logger, sentry, monitor, and apps extensions
- grpc - Blank application with grpc extension (optional client)
- http - Blank application with http extension
