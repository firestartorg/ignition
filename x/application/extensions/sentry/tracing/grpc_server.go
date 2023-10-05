package tracing

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		// Add a sentry hub to the request context
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			log.Ctx(ctx).Warn().Msg("No sentry hub found in context")
			return handler(ctx, req)
		}

		options := []sentry.SpanOption{
			sentry.WithOpName("grpc.server"),
			ContinueFromRpcRequest(ctx),
			sentry.WithTransactionSource(sentry.SourceURL),
		}

		// Start a new transaction
		transaction := sentry.StartTransaction(ctx,
			info.FullMethod,
			options...,
		)
		defer transaction.Finish()

		// TODO: Perhaps makes sense to use SetRequestBody instead?
		// Add the request to the scope
		hub.Scope().SetExtra("requestBody", req)
		// Add the tags from grpc request metadata
		tags := grpc_ctxtags.Extract(ctx)
		for k, v := range tags.Values() {
			hub.Scope().SetTag(k, v.(string))
		}

		resp, err := handler(ctx, req)
		if err != nil {
			hub.CaptureException(err)
		}

		// Add the response status to the scope
		transaction.Status = ToSpanStatus(status.Code(err))

		return resp, err
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {

		ctx := ss.Context()

		// Add a sentry hub to the request context
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			log.Ctx(ctx).Warn().Msg("No sentry hub found in context")
			return handler(ctx, ss)
		}

		options := []sentry.SpanOption{
			sentry.WithOpName("grpc.server"),
			ContinueFromRpcRequest(ctx),
			sentry.WithTransactionSource(sentry.SourceURL),
		}

		// Start a new transaction
		transaction := sentry.StartTransaction(ctx,
			info.FullMethod,
			options...,
		)
		defer transaction.Finish()

		// Wrap the server stream with a new context
		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = ctx
		ss = wrapped

		// Add the tags from grpc request metadata
		tags := grpc_ctxtags.Extract(ctx)
		for k, v := range tags.Values() {
			hub.Scope().SetTag(k, v.(string))
		}

		err := handler(ctx, ss)
		if err != nil {
			hub.CaptureException(err)
		}

		// Add the response status to the scope
		transaction.Status = ToSpanStatus(status.Code(err))

		return err
	}
}
