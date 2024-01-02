# HTTP Extension

The HTTP extension is used to add one or more HTTP servers to your application.

## How to add the extension

The extension needs to be configured in code.

```go
package main

import (
  "gitlab.com/firestart/ignition/x/application"
  "gitlab.com/firestart/ignition/x/application/extensions/http"
)

func main() {
  app := application.New(
    http.WithServer(
      http.WithPort(3000), // default port is 3000
      http.WithMiddleware(/* your middleware functions */)
    ),
    // You can add as many servers as you want, but they need to have different ports
    // and you need to give them a name
    http.WithNamedServer(
      "admin",
      http.WithPort(3001),
      http.WithMiddleware(/* your middleware functions */)
    ),
  )

  app.Run()
}
```

## How to add routes

Routes can be added to the HTTP extension using hooks.

```go
package main

import (
  "gitlab.com/firestart/ignition/x/application"
  "gitlab.com/firestart/ignition/x/application/extensions/http"
  nethttp "net/http"
  "github.com/julienschmidt/httprouter"
)

func main() {
  app := application.New(
    http.WithServer())

  // You can add routes to the default server
  http.MustAddGetRoute("/", func(w nethttp.ResponseWriter, r *nethttp.Request, _ httprouter.Params) {
    w.Write([]byte("Hello World"))
  })

  // If you want to add a route to a different server, you need to use the following code
  // http.AddNamedRoute
  // http.MustAddNamedRoute

  app.Run()
}
```

## How to add request context processors

Request context processors can be used for multiple things, like:

- adding a logger to the request context
- adding a request ID to the request context
- passing tracing information to http clients

To add a request context processor, you need to use the following code:

```go
...

func main() {
  app := application.New(
    http.WithServer())

  app.AddContextProcessor(application.HookRequest, func(ctx context.Context, app App) (context.Context, error) {
    // Add your code here
    return ctx, nil
  })

  app.Run()
}
```

There are extensions that add request context processors, like:

- [sentry](../sentry/README.md)
- logging
