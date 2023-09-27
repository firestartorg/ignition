package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/firestart/ignition/examples/application/pb"
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/apps"
	"gitlab.com/firestart/ignition/x/application/extensions/grpc"
	"gitlab.com/firestart/ignition/x/application/extensions/http"
	"gitlab.com/firestart/ignition/x/application/extensions/monitor"
	http1 "net/http"
)

func main() {
	// Create the application
	app := application.New(
		apps.WithVersion("sample", "1.0.0"),
		apps.WithHealthCheck(func(ctx context.Context, app application.App) error {
			return nil
		}),
		monitor.WithDefaultMonitor(),

		grpc.WithClientFactory(),
		grpc.WithServer(),

		http.WithServer(),
	)

	// Setup the gRPC server
	grpc.MustUseReflection(app)
	grpc.MustAddService(app, pb.Greeter_ServiceDesc, Greeter{})

	// Setup the HTTP server
	http.MustAddGetRoute(app, "/hello", func(w http1.ResponseWriter, r *http1.Request, ps httprouter.Params) {
		_, _ = w.Write([]byte("Hello World"))
	})

	// Run the application
	app.Run()
}

type Greeter struct{}

func (g Greeter) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + request.Name}, nil
}

var _ pb.GreeterServer = (*Greeter)(nil)
