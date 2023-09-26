package monitor

import "gitlab.com/firestart/ignition/x/application"

var (
	HookHealth    application.Hook = "health"
	HookReadiness application.Hook = "health/ready"
	HookLiveness  application.Hook = "health/live"
)
