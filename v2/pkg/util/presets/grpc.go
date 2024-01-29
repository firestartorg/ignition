package presets

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/grpc"
	"gitlab.com/firestart/ignition/x/application/extensions/sentry/recovery"
	"gitlab.com/firestart/ignition/x/application/extensions/sentry/tracing"
	"gitlab.com/firestart/ignition/x/goenv"
	grpc1 "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// WithRpcClientFactory adds the grpc client to the application
func WithRpcClientFactory() application.Option {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// AddHook any other option (check functions starting with logging.With).
	}

	return grpc.WithClientFactory(
		grpc1.WithTransportCredentials(insecure.NewCredentials()),
		grpc1.WithChainUnaryInterceptor(
			tracing.UnaryClientInterceptor(),
			logging.UnaryClientInterceptor(GrpcLogger, opts...),
		),
		grpc1.WithChainStreamInterceptor(
			tracing.StreamClientInterceptor(),
			logging.StreamClientInterceptor(GrpcLogger, opts...),
		),
	)
}

// WithRpcServer adds the grpc server to the application
func WithRpcServer(port int16) application.Option {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// AddHook any other option (check functions starting with logging.With).
	}

	return grpc.WithServerPort(
		port,
		grpc1.ChainUnaryInterceptor(
			tracing.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recovery.RpcRecoveryHandler)),
			logging.UnaryServerInterceptor(GrpcLogger, opts...),
		),
		grpc1.ChainStreamInterceptor(
			// tracing.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recovery.RpcRecoveryHandler)),
			logging.StreamServerInterceptor(GrpcLogger, opts...),
		),
	)
}

// NewRpcApp creates a new application with the following components:
// - blank app preset (see NewBlankApp)
// - grpc server (see WithRpcServer)
// - grpc reflection (if not in production)
func NewRpcApp(name string, opts ...application.Option) application.App {
	app := NewBlankApp(
		name,
		pack(
			opts,
			[]application.Option{
				MakeConfigurable("App", func(config grpcConfig) []application.Option {
					// SetInjectable the default port
					if config.Port == 0 {
						config.Port = 5000
					}

					return []application.Option{
						WithRpcServer(config.Port),
					}
				}),
			},
		)...,
	)

	if !goenv.IsProduction() {
		grpc.MustUseReflection(app)
	}

	return app
}

type grpcConfig struct {
	// Port is the port to listen on
	Port int16
}
