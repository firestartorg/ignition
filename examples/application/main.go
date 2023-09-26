package main

import (
	"context"
	"gitlab.com/firestart/ignition/examples/application/pb"
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/grpc"
	"gitlab.com/firestart/ignition/x/application/monitor"
)

func main() {
	// Create the application
	app := application.New(
		monitor.WithAppVersion("sample", "1.0.0"),
		monitor.WithDefaultMonitor(),

		grpc.WithClientFactory(),
		grpc.WithServer(),
	)

	// Inject dependencies
	grpc.MustUseReflection(app)
	grpc.MustAddService(app, pb.Greeter_ServiceDesc, Greeter{})

	// Run the application
	app.Run()
}

type Greeter struct{}

func (g Greeter) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + request.Name}, nil
}

var _ pb.GreeterServer = (*Greeter)(nil)
