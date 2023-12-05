package http

import (
	"fmt"
	"net/http"
)

type ServerOption = func(opts *ServerOptions)
type ServerOptions struct {
	// port is the port to listen on
	port int16
	// middleware is a list of middleware to apply to the server
	middlewares []Middleware
}

// newServerOptions creates a new server options struct
func newServerOptions(options ...ServerOption) ServerOptions {
	opts := ServerOptions{
		port:        3000,
		middlewares: []Middleware{},
	}

	for _, option := range options {
		option(&opts)
	}

	return opts
}

// addr returns the address to listen on
func (opts *ServerOptions) addr() string {
	return fmt.Sprint(":", opts.port)
}

// wrapMiddleware wraps the handler with the middleware
func (opts *ServerOptions) wrapMiddleware(handler http.Handler) http.Handler {
	for i := len(opts.middlewares) - 1; i >= 0; i-- {
		handler = opts.middlewares[i](handler)
	}
	return handler
}

// WithPort sets the port to listen on
func WithPort(port int16) ServerOption {
	return func(opts *ServerOptions) {
		opts.port = port
	}
}

// WithMiddleware adds middleware to the server
func WithMiddleware(middleware ...Middleware) ServerOption {
	return func(opts *ServerOptions) {
		opts.middlewares = append(opts.middlewares, middleware...)
	}
}
