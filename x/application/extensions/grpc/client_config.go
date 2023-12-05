package grpc

import (
	"fmt"
	"gitlab.com/firestart/ignition/pkg/injector"
)

// NewConfiguredClient creates a new client from the application configuration.
// The configuration should be in the following format:
//
//	grpc:
//	  client:
//	    <name>:
//	      host: <host>
func NewConfiguredClient[T interface{}](inj *injector.Injector, name string, f NewClientFunc[T]) (T, error) {
	// Get configuration for the client
	conf, err := injector.ExtractConfig[clientConfig](inj, fmt.Sprint("Grpc:Client:", name))
	if err != nil {
		var zero T
		return zero, err
	}

	// Create the client
	return NewClient(inj, conf.Host, f)
}

// MustNewConfiguredClient creates a new client and panics if there is an error
func MustNewConfiguredClient[T interface{}](inj *injector.Injector, name string, f NewClientFunc[T]) T {
	// Create the client
	c, err := NewConfiguredClient[T](inj, name, f)
	if err != nil {
		panic(err)
	}
	return c
}

type clientConfig struct {
	// Host is the host to connect to
	Host string
}
