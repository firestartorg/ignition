package presets

import (
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/grpc"
	"gitlab.com/firestart/ignition/x/goenv"
)

// NewRpcApp creates a new application with the following components:
// - blank app preset (see NewBlankApp)
// - grpc server
// - grpc reflection (if not in production)
func NewRpcApp(name string, opts ...application.Option) application.App {
	app := NewBlankApp(
		name,
		pack(
			opts,
			[]application.Option{
				grpc.WithServer(),
			},
		)...,
	)

	if !goenv.IsProduction() {
		grpc.MustUseReflection(app)
	}

	return app
}
