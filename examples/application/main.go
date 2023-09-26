package main

import (
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/grpc"
)

func main() {
	// Create the application
	app := application.New(
		grpc.WithClientFactory(),
		grpc.WithServer(),
	)

	// Inject dependencies
	grpc.MustAddService(app, nil, nil)

	// Run the application
	app.Run()
}
