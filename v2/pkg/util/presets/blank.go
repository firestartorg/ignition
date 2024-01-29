package presets

import (
	"github.com/rs/zerolog"
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/config"
	"gitlab.com/firestart/ignition/x/application/extensions/logging"
	"gitlab.com/firestart/ignition/x/application/extensions/monitor"
	"gitlab.com/firestart/ignition/x/application/extensions/sentry"
	"gitlab.com/firestart/ignition/x/application/extensions/vcs"
	"os"
)

// NewBlankApp creates a new application with the following extensions:
// - vcs.WithBuildInfo
// - config.WithConfiguration
// - sentry.WithDefaultSentry
// - monitor.WithDefaultMonitor
// - logging.WithZerolog
func NewBlankApp(name string, opts ...application.Option) application.App {
	logger := zerolog.New(os.Stderr).With().Timestamp().Stack().Caller().Logger()

	return application.New(pack(
		[]application.Option{
			vcs.WithBuildInfo(name),
			config.WithConfiguration(),
			sentry.WithDefaultSentry(),
			monitor.WithMonitor(
				monitor.FromConfig(),
				monitor.WithPrometheusMetrics(),
				monitor.WithReadiness(),
				monitor.WithLiveness(),
			),
			logging.WithConfigurableZerolog(logger),
		},
		opts,
	)...)
}
