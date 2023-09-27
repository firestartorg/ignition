package http

import "fmt"

type ServerOptions struct {
	port int16
}

type ServerOption = func(opts *ServerOptions)

func (opts *ServerOptions) addr() string {
	return fmt.Sprint(":", opts.port)
}
