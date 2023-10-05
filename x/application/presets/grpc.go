package presets

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
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
	return grpc.WithClientFactory(
		grpc1.WithTransportCredentials(insecure.NewCredentials()),
		grpc1.WithChainUnaryInterceptor(tracing.UnaryClientInterceptor()),
		grpc1.WithChainStreamInterceptor(tracing.StreamClientInterceptor()),
	)
}

// WithRpcServer adds the grpc server to the application
func WithRpcServer(port int16) application.Option {
	return grpc.WithServerPort(
		port,
		grpc1.ChainUnaryInterceptor(
			tracing.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recovery.RpcRecoveryHandler)),
		),
		grpc1.ChainStreamInterceptor(
			tracing.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recovery.RpcRecoveryHandler)),
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
					// Set the default port
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
