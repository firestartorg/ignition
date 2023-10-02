package presets

import (
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/config"
	"gitlab.com/firestart/ignition/x/application/extensions/grpc"
	"gitlab.com/firestart/ignition/x/application/extensions/logging"
	"gitlab.com/firestart/ignition/x/application/extensions/monitor"
	"gitlab.com/firestart/ignition/x/application/extensions/sentry"
	"gitlab.com/firestart/ignition/x/goenv"
)

func NewRpcApp(opts ...application.Option) application.App {
	app := application.New(pack(
		[]application.Option{
			config.WithSettings(),
			sentry.WithDefaultSentry(),
			monitor.WithDefaultMonitor(),
			logging.WithZerolog(),
		},
		opts,
		[]application.Option{
			grpc.WithServer(),
		},
	)...)

	if !goenv.IsProduction() {
		grpc.MustUseReflection(app)
	}

	return app
}
