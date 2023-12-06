package grpc

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
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

// WithServerPort sets the server
func WithServerPort(port int16, opts ...grpc.ServerOption) application.Option {
	return func(app application.App) {
		// Pack the server options
		opts = packServer(
			[]grpc.ServerOption{
				grpc.ChainUnaryInterceptor(UnaryServerInterceptor(app)),
				grpc.ChainStreamInterceptor(StreamServerInterceptor(app)),
			},
			opts)
		// Create the server
		srv := grpc.NewServer(opts...)

		// AddHook the server to the application
		injector.InjectNamed(app.Injector, ServerName, server{
			server: srv,
		})

		// AddHook the startup hook
		app.AddStartupHook(func(ctx context.Context, app application.App) error {
			//fmt.Println("Starting gRPC server")

			log.Ctx(ctx).Info().
				Int16("port", port).
				Msg("Starting gRPC server")

			// Start the server
			listen, err := net.Listen("tcp", fmt.Sprint(":", port))
			if err != nil {
				return err
			}
			return srv.Serve(listen)
		})

		// AddHook the shutdown hook
		app.AddShutdownHook(func(ctx context.Context, app application.App) error {
			// Stop the server
			srv.GracefulStop()

			return nil
		})
	}
}

// WithServer adds a server to the application
func WithServer(opts ...grpc.ServerOption) application.Option {
	return WithServerPort(5000, opts...)
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
