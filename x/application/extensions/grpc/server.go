package grpc

import (
	"context"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

var (
	// ServerName is the name of the injectable that contains the server
	ServerName = "ignition/grpc-server"
)

// server is a container for the server
type server struct {
	// server is the server
	server *grpc.Server
}

// WithServer sets the server
func WithServer(opts ...grpc.ServerOption) application.Option {
	return func(app application.App, hooks *application.Hooks) {
		// Pack the server options
		opts = packServer(
			[]grpc.ServerOption{
				grpc.ChainUnaryInterceptor(UnaryServerInterceptor(app, hooks)),
				grpc.ChainStreamInterceptor(StreamServerInterceptor(app, hooks)),
			},
			opts)
		// Create the server
		srv := grpc.NewServer(opts...)

		// Add the server to the application
		injector.InjectNamed(app.Injector, ServerName, server{
			server: srv,
		})

		// Add the startup hook
		hooks.AddStartup(func(ctx context.Context, app application.App) error {
			//fmt.Println("Starting gRPC server")

			// Start the server
			listen, err := net.Listen("tcp", ":5000")
			if err != nil {
				return err
			}
			return srv.Serve(listen)
		})

		// Add the shutdown hook
		hooks.AddShutdown(func(ctx context.Context, app application.App) error {
			// Stop the server
			srv.GracefulStop()

			return nil
		})
	}
}

func AddService(app application.App, desc grpc.ServiceDesc, srv any) error {
	// Get the server
	s, err := injector.GetNamed[server](app.Injector, ServerName)
	if err != nil {
		return err
	}

	// Register the service
	s.server.RegisterService(&desc, srv)

	return nil
}

// MustAddService adds a service to the application, panicking if it fails
func MustAddService(app application.App, desc grpc.ServiceDesc, srv any) {
	if err := AddService(app, desc, srv); err != nil {
		panic(err)
	}
}

// MustUseReflection adds reflection to the grpc application.
// Experimental, may be removed, changed or replaced with a better solution in the future.
func MustUseReflection(app application.App) {
	s, err := injector.GetNamed[server](app.Injector, ServerName)
	if err != nil {
		panic(err)
	}

	reflection.Register(s.server)
}
