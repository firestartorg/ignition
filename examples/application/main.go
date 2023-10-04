package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/examples/application/pb"
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/apps"
	"gitlab.com/firestart/ignition/x/application/extensions/grpc"
	"gitlab.com/firestart/ignition/x/application/extensions/http"
	"gitlab.com/firestart/ignition/x/application/extensions/logging"
	"gitlab.com/firestart/ignition/x/application/extensions/monitor"
	"gitlab.com/firestart/ignition/x/application/extensions/vcs"
	http1 "net/http"
)

func main() {
	// Create the application
	app := application.New(
		vcs.WithBuildInfo("examples/application"),
		apps.WithHealthCheck(func(ctx context.Context, app application.App) error {
			return nil
		}),
		monitor.WithDefaultMonitor(),

		logging.WithDefaultZerolog(),

		grpc.WithClientFactory(),
		grpc.WithServer(),

		http.WithServer(),
	)

	// Setup the gRPC server
	grpc.MustUseReflection(app)
	grpc.MustAddService(app, pb.Greeter_ServiceDesc, Greeter{})

	// Setup the HTTP server
	http.MustAddGetRoute(app, "/hello", func(w http1.ResponseWriter, r *http1.Request, ps httprouter.Params) {
		log.Ctx(r.Context()).Info().Msg("Hello World")

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
