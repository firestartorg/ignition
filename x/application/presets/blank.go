package presets

import (
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/config"
	"gitlab.com/firestart/ignition/x/application/extensions/logging"
	"gitlab.com/firestart/ignition/x/application/extensions/monitor"
	"gitlab.com/firestart/ignition/x/application/extensions/sentry"
	"gitlab.com/firestart/ignition/x/application/extensions/vcs"
)

// NewBlankApp creates a new application with the following extensions:
// - vcs.WithBuildInfo
// - config.WithSettings
// - sentry.WithDefaultSentry
// - monitor.WithDefaultMonitor
// - logging.WithZerolog
func NewBlankApp(name string, opts ...application.Option) application.App {
	return application.New(pack(
		[]application.Option{
			vcs.WithBuildInfo(name),
			config.WithSettings(),
			sentry.WithDefaultSentry(),
			monitor.WithDefaultMonitor(),
			logging.WithDefaultZerolog(),
		},
		opts,
	)...)
}
