package main

import (
	"context"
	sentry1 "github.com/getsentry/sentry-go"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/examples/application/pb"
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/apps"
	"gitlab.com/firestart/ignition/x/application/extensions/config"
	"gitlab.com/firestart/ignition/x/application/extensions/grpc"
	"gitlab.com/firestart/ignition/x/application/extensions/http"
	"gitlab.com/firestart/ignition/x/application/extensions/logging"
	"gitlab.com/firestart/ignition/x/application/extensions/monitor"
	"gitlab.com/firestart/ignition/x/application/extensions/sentry"
	"gitlab.com/firestart/ignition/x/application/extensions/vcs"
	"gitlab.com/firestart/ignition/x/application/presets"
	http1 "net/http"
	"os"
)

func main() {
	// Create the application
	app := application.New(
		vcs.WithBuildInfo("examples/application"),
		apps.WithHealthCheck(func(ctx context.Context, app application.App) error {
			return nil
		}),
		config.WithSettings(),
		monitor.WithMonitor(
			monitor.FromConfig(),
			monitor.WithPrometheusMetrics(),
			monitor.WithReadiness(),
			monitor.WithLiveness(),
		),

		sentry.WithSentry(
			sentry.WithDsn("https://912642dbd878f3999a37d9d42937a5ec@o4505345231945728.ingest.sentry.io/4505992844804096"),
			sentry.EnableTracing(),
		),
		logging.WithZerolog(
			zerolog.New(os.Stderr).With().
				Timestamp().
				Stack().
				Caller().
				Logger().
				Hook(sentry.LoggerHook{}),
		),

		presets.WithRpcClientFactory(),
		presets.WithRpcServer(5001),
		presets.WithHttpServer(3000),
	)

	// Setup the gRPC server
	grpc.MustUseReflection(app)
	ProvideGreeterService(app)

	// Setup the HTTP server
	http.MustAddGetRoute(app, "/hello", func(w http1.ResponseWriter, r *http1.Request, ps httprouter.Params) {
		ctx := r.Context()

		// Create a new span
		span := sentry1.StartSpan(ctx, "function")
		span.Description = "suboperation2"

		// Do some work
		log.Ctx(ctx).Info().Msg("Hello World")
		_, _ = w.Write([]byte("Hello World"))

		// Finish the span
		span.Finish()
	})
	http.MustAddGetRoute(app, "/panic", func(w http1.ResponseWriter, r *http1.Request, ps httprouter.Params) {
		panic("Panic!")
	})

	// Run the application
	app.Run()
}

type Greeter struct {
	client pb.GreeterClient
}

func (g Greeter) Panic(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	panic(request.Name)
}

func (g Greeter) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Ctx(ctx).Info().Msgf("Received: %s", request.Name)

	_, err := g.client.Panic(ctx, &pb.HelloRequest{Name: "Panic!"})
	if err != nil {
		return nil, err
	}

	return &pb.HelloReply{Message: "Hello " + request.Name}, nil
}

var _ pb.GreeterServer = (*Greeter)(nil)

func ProvideGreeterService(app application.App) {
	srv := Greeter{
		client: grpc.MustNewClient(app.Injector, "localhost:5000", pb.NewGreeterClient),
	}
	grpc.MustAddService(app, pb.Greeter_ServiceDesc, srv)
	return
}
