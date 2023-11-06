package grpc

import (
	"context"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/monitor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync"
)

var (
	// ClientFactoryName is the name of the injectable that contains the client clientFactory
	ClientFactoryName = "ignition/grpc-client-clientFactory"
)

// clientFactory is a container for the client clientFactory
type clientFactory struct {
	// clients is a list of clients
	clients []*grpc.ClientConn
	// clientsMutex is a mutex to protect the clients list
	clientsMutex sync.RWMutex

	// new generates a new client
	// and stores it in the clients list
	new func(target string) (*grpc.ClientConn, error)
}

// WithClientFactory adds a grpc client factory to the application.
// Features:
//   - automatically monitors all clients
//   - automatically closes all clients
func WithClientFactory(opts ...grpc.DialOption) application.Option {
	return func(app application.App) {
		// Create the client factory
		factory := &clientFactory{clients: make([]*grpc.ClientConn, 0), clientsMutex: sync.RWMutex{}}
		factory.new = func(target string) (*grpc.ClientConn, error) {
			conn, err := grpc.Dial(target, opts...)
			if err != nil {
				return nil, err
			}

			// Add the client to the list
			factory.clientsMutex.Lock()
			defer factory.clientsMutex.Unlock()
			factory.clients = append(factory.clients, conn)

			return conn, nil
		}

		// Add the client clientFactory to the application
		injector.InjectNamed(app.Injector, ClientFactoryName, factory)

		// Add the shutdown hook
		app.AddShutdownHook(func(ctx context.Context, app application.App) error {
			// Close all the clients
			factory.clientsMutex.Lock()
			defer factory.clientsMutex.Unlock()

			for _, client := range factory.clients {
				if client.GetState() == connectivity.Shutdown {
					continue
				}

				err := client.Close()
				if err != nil {
					return err
				}
			}
			return nil
		})

		// Add health check hook
		app.AddHook(monitor.HookReadiness, func(ctx context.Context, app application.App) error {
			// Check all the clients
			factory.clientsMutex.RLock()
			defer factory.clientsMutex.RUnlock()

			for _, client := range factory.clients {
				if client.GetState() != connectivity.Ready {
					log.Ctx(ctx).Error().Str("target", client.Target()).Msg("grpc client is not ready")
					return monitor.ErrNotReady
				}
			}
			return nil
		})
	}
}

// NewClientConnection creates a new client
func NewClientConnection(inj *injector.Injector, target string) (*grpc.ClientConn, error) {
	// Get the client clientFactory
	f, err := injector.GetNamed[*clientFactory](inj, ClientFactoryName)
	if err != nil {
		return nil, err
	}

	return f.new(target)
}

// MustNewClientConnection creates a new client and panics if there is an error
func MustNewClientConnection(inj *injector.Injector, target string) *grpc.ClientConn {
	// Create the client
	conn, err := NewClientConnection(inj, target)
	if err != nil {
		panic(err)
	}
	return conn
}

type NewClientFunc[T interface{}] func(conn grpc.ClientConnInterface) T

// NewClient creates a new client
func NewClient[T interface{}](inj *injector.Injector, target string, f NewClientFunc[T]) (c T, err error) {
	// Create the client
	var conn *grpc.ClientConn
	conn, err = NewClientConnection(inj, target)
	if err != nil {
		return
	}

	// Create the client
	c = f(conn)
	return
}

// MustNewClient creates a new client and panics if there is an error
func MustNewClient[T interface{}](inj *injector.Injector, target string, f NewClientFunc[T]) T {
	// Create the client
	c, err := NewClient[T](inj, target, f)
	if err != nil {
		panic(err)
	}
	return c
}
