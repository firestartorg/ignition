package vcs

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/x/application"
	"runtime/debug"
)

var (
	// BuildInfo is the gauge metric exposed
	BuildInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "build_info",
			Help: "Metric with a constant '1' value labeled by version and goversion from which the cloud gateway receiver was built.",
		},
		[]string{"app_name", "app_version", "go_version"},
	)
)

// WithBuildInfo set the application build info (version and go version) metric
func WithBuildInfo(name string) application.Option {
	return func(app application.App, hooks *application.Hooks) {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Info().Msg("No build info found")
		}

		revision, ok := findSetting(info.Settings, "vcs.revision")
		if !ok {
			// Add the following metric
			BuildInfo.WithLabelValues(name, "dev", info.GoVersion).Set(1)
			return
		}

		BuildInfo.WithLabelValues(name, revision.Value, info.GoVersion).Set(1)
	}
}

func findSetting(settings []debug.BuildSetting, key string) (debug.BuildSetting, bool) {
	for _, setting := range settings {
		if setting.Key == key {
			return setting, true
		}
	}

	return debug.BuildSetting{}, false
}
