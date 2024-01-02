# Monitor Extension

The Monitor extension provides a simple way to monitor the status of your application using:

- a metrics endpoint `/metrics`
- a health endpoint
  - a readiness endpoint `/health/ready`
  - a liveness endpoint `/health/live`

## How to add the extension

The extension can be configured in code.

```go
package main

import (
  "gitlab.com/firestart/ignition/x/application"
  "gitlab.com/firestart/ignition/x/application/extensions/monitor"
)

func main() {
  app := application.New(
    monitor.WithMonitor(
      monitor.WithPrometheusMetrics(),
      monitor.WithReadiness(),
      monitor.WithLiveness(),
      // If you want to change the port
      // You can use the WithPort option
      monitor.WithPort(8080),
      // Or by loading the configuration from a file
      // You can use the FromConfig option (App:Monitor in config)
      monitor.FromConfig(),
      // Or by automatically assigning a port
      // You can use the WithPortAntiCollision option
      monitor.WithPortAntiCollision(),
    ),
  )
  app.Run()
}
```

## How to add health checks

Health checks can be added to the monitor extension using hooks.

```go
package main

import (
  "gitlab.com/firestart/ignition/x/application"
  "gitlab.com/firestart/ignition/x/application/extensions/monitor"
  "context"
)

func main() {
  app := application.New(
    monitor.WithMonitor(
      monitor.WithReadiness(),
      monitor.WithLiveness(),
    ),
  )
  // Hook to add the liveness and readiness checks
  app.AddHook(monitor.HookHealth, func(ctx context.Context, app application.App) error {
    return nil // Return an error to make the application unhealthy
  })
  // Hook to add the liveness checks
  app.AddHook(monitor.HookLiveness, func(ctx context.Context, app application.App) error {
    return nil // Return an error to make the application unhealthy
  })
  // Hook to add the readiness checks
  app.AddHook(monitor.HookReadiness, func(ctx context.Context, app application.App) error {
    return nil // Return an error to make the application not ready
  })

  app.Run()
}
```

Some extensions automatically add health checks:

- [gRPC](../grpc/README.md): adds a health check for the gRPC server and clients
- mongo: adds a health check for the mongo client
