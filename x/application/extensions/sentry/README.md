# Sentry Extension

The Sentry extension is used to add Sentry error reporting to your application. It uses 
the [sentry-go](https://sentry.io/for/go/) library. 

## How to add the extension

The extension needs to be configured in code, but can optionally be configured using the 
`config` extension.

```go
package main

import (
  "gitlab.com/firestart/ignition/x/application"
  "gitlab.com/firestart/ignition/x/application/extensions/sentry"
)

func main() {
  app := application.New(
    sentry.WithSentry(
      sentry.WithDsn("https://example.com"),
      sentry.WithRelease("1.0.0"),
      sentry.WithEnvironment("production"),
      // Or you can use the config extension
      sentry.FromSettings()
    )
  )

  app.Run()
}
```

The configuration options are:

- `enabled`: Whether to enable Sentry or not
- `dsn`: The DSN of your Sentry project
- `release`: The release version of your application
- `environment`: The environment of your application
- `debug`: Whether to enable debug mode or not
- `enableTracing`: Whether to enable tracing or not
- `tracingSampleRate`: The tracing sample rate
- `flushTimeout`: The timeout for flushing events to Sentry

A sample configuration file can be found in the [examples](../../../../examples/application/rpc/appsettings.yaml) directory.
