package grpc

import (
	"context"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/x/application"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor is a grpc server interceptor that processes the request context
func UnaryServerInterceptor(app application.App, hooks *application.Hooks) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		var err error
		// Process the request context
		ctx, err = hooks.ProcessContext(application.HookRequest, ctx, app)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to process request context")
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor is a grpc server interceptor that processes the request context
func StreamServerInterceptor(app application.App, hooks *application.Hooks) grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		var err error
		// Process the request context
		ctx, err := hooks.ProcessContext(application.HookRequest, stream.Context(), app)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to process request context")
		} else {
			wrapped := grpcmiddleware.WrapServerStream(stream)
			wrapped.WrappedContext = ctx

			stream = wrapped
		}

		return handler(srv, stream)
	}
}

// packServer packs the server options
func packServer(opts ...[]grpc.ServerOption) []grpc.ServerOption {
	var packed []grpc.ServerOption
	for _, o := range opts {
		if o == nil {
			continue
		}
		packed = append(packed, o...)
	}
	return packed
}
