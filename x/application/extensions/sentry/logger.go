package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

// TODO: Replace with custom Writer that uses the sentry client

type LoggerHook struct {
}

func (h LoggerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.PanicLevel && level != zerolog.FatalLevel && level != zerolog.ErrorLevel {
		return
	}

	// Get the context from the event
	ctx := e.GetCtx()
	// Get the hub from the context
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	// Capture the event
	hub.CaptureMessage(msg)
}
