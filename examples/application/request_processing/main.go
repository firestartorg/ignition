package main

import (
	"context"
	"github.com/google/uuid"
	"gitlab.com/firestart/ignition/x/application"
)

func main() {
	// This is a simple example of how preprocessor each request context.

	// First we need to create a new application.
	// One could use a preset, like so:
	// app := presets.NewHttpApp("examples/application/request_processing")
	// But for this example we will create a blank application.
	app := application.New()

	// Then we can register a preprocessor.
	// A preprocessor is a function that is called before each request.
	app.AddContextProcessor(application.HookRequest, func(ctx context.Context, app application.App) (context.Context, error) {
		// This is the function that will be called before each request.
		return ctx, nil
	})

	// For example, we could add a unique id to each request.
	app.AddContextProcessor(application.HookRequest, func(ctx context.Context, app application.App) (context.Context, error) {
		requestId, _ := uuid.NewRandom()

		return context.WithValue(ctx, "requestId", requestId), nil
	})
}
