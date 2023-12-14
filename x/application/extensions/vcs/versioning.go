package vcs

import (
	"runtime/debug"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/x/application"
)

var (
	Branch string = "undefined"
	// BuildInfo is the gauge metric exposed
	BuildInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "build_info",
			Help: "Metric with a constant '1' value labeled by app name, app version, git branch and goversion from which the application was built.",
		},
		[]string{"app_name", "app_version", "git_branch", "go_version"},
	)
)

// WithBuildInfo set the application build info (version and go version) metric
func WithBuildInfo(name string) application.Option {
	return func(app application.App) {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Info().Msg("No build info found")
		}

		revision, ok := findSetting(info.Settings, "vcs.revision")
		if !ok {
			// AddHook the following metric
			BuildInfo.WithLabelValues(name, "dev", Branch, info.GoVersion).Set(1)
			return
		}

		BuildInfo.WithLabelValues(name, revision.Value, Branch, info.GoVersion).Set(1)
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
