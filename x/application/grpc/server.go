package grpc

import (
	"context"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
	"google.golang.org/grpc"
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
		// Create the server
		srv := grpc.NewServer(opts...)

		// Add the server to the application
		injector.InjectNamed(app.Injector, ServerName, server{
			server: srv,
		})

		// Add the startup hook
		hooks.AddStartup(func(ctx context.Context, app application.App) error {
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

type RegisterFunc = func(s grpc.ServiceRegistrar, srv any)

func AddService(app application.App, register RegisterFunc, srv any) error {
	// Get the server
	s, err := injector.GetNamed[server](app.Injector, ServerName)
	if err != nil {
		return err
	}

	// Register the service
	register(s.server, srv)

	return nil
}

// MustAddService adds a service to the application, panicking if it fails
func MustAddService(app application.App, register RegisterFunc, srv any) {
	if err := AddService(app, register, srv); err != nil {
		panic(err)
	}
}
