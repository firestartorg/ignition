# gRPC Extension

The gRPC extension can be used to

- add a gRPC [server](#server) to your application
- add a gRPC [client](#client) to your application
- add a gRPC [client factory](#client) to your application

## Server

In order to add a gRPC server to your application, you need to use the `grpc.WithServer` option.

```go
package main

import (
  "gitlab.com/firestart/ignition/x/application"
  "gitlab.com/firestart/ignition/x/application/extensions/grpc"
)

func main() {
  app := application.New(
    grpc.WithServer(
      // Here you can add normal gRPC server options, like grpc.WithChainUnaryInterceptors
      // Note that the server as it is is already configured with a Chain stream and unary interceptors
      // so dont use ChainUnaryInterceptor or ChainStreamInterceptor here unless you want to override the default ones
      // If you want to add additional interceptors, use grpc.WithUnaryInterceptors and grpc.WithStreamInterceptors
    ),
  )

  // You can add services to the server using the following code
  grpc.MustAddService(
    desc, // This is the service descriptor, you can get it by importing the proto file (const desc = pb.<Service>_ServiceDesc)
    &Service{},
  )

  app.Run()
}
```

See the example in the [examples/application/rpc](../../../../examples/application/rpc) folder for a full example.

## Client

In order to add a gRPC client to your application, you need to use the `grpc.WithClientFactory` option.

```go
package main

import (
  "gitlab.com/firestart/ignition/x/application"
  "gitlab.com/firestart/ignition/x/application/extensions/grpc"
)

func main() {
  app := application.New(
    grpc.WithClientFactory(
      // Here you can add normal gRPC client options, like grpc.WithChainUnaryInterceptors
      // Note that the client as it is is already configured with a Chain stream and unary interceptors
      // so dont use ChainUnaryInterceptor or ChainStreamInterceptor here unless you want to override the default ones
      // If you want to add additional interceptors, use grpc.WithUnaryInterceptors and grpc.WithStreamInterceptors
    ),
  )

  app.Run()
}
```

### Client Connection

You can create a new client connection using the `grpc.NewClientConnection` function.

```go
conn, dial := grpc.NewClientConnection(
  injector, // This can be app.Injector or a injector reference provided in a Provide function
  "localhost:3000", // This is the address of the server
) // Or MustNewClientConnection if you want to panic on error
```

### Client

You can create a new client using the `grpc.NewClient` function.

```go
client, err := grpc.NewClient(
  injector, // This can be app.Injector or a injector reference provided in a Provide function
  "localhost:3000", // This is the address of the server
  pb.NewServiceClient, // This the function provided by the proto file (pb.New<Service>Client)
) // Or MustNewClient if you want to panic on error
```

### Configured Client

You can create a new client using the `grpc.NewConfiguredClient` function. For this to work, 
you need to add the `config` extension to your application.

```go
client, err := grpc.NewConfiguredClient(
  injector, // This can be app.Injector or a injector reference provided in a Provide function
  "Greeter", // This is the name of the client in the config file
  //            This name needs to start with a capital letter (see config package for more info)
  pb.NewServiceClient, // This the function provided by the proto file (pb.New<Service>Client)
) // Or MustNewConfiguredClient if you want to panic on error
```

The application config file needs to have the following structure:

```yaml
grpc:
  client:
    greeter: # This name needs to match the name provided in the NewConfiguredClient function
             # The name needs to start with a lowercase letter (see config package for more info)
      host: localhost:3000
```