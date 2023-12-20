package tracing

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		callOpts ...grpc.CallOption,
	) error {
		// Create a new span
		span := sentry.StartSpan(ctx, "grpc.client")
		span.Description = method
		// If finished, the span will be sent to sentry
		defer span.Finish()

		// AddHook the sentry trace and baggage headers to the metadata
		ctx = getSpanContext(span)

		err := invoker(ctx, method, req, reply, cc, callOpts...)
		if err != nil && !IsGrpcError(err) {
			hub := sentry.GetHubFromContext(ctx)
			if hub != nil {
				hub.CaptureException(err)
			} else {
				log.Ctx(ctx).Warn().Msg("No sentry hub found in context")
			}
		}

		// AddHook the response status to the span
		span.Status = ToSpanStatus(status.Code(err))

		return err
	}
}

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		callOpts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		// Create a new span
		span := sentry.StartSpan(ctx, "grpc.client")
		span.Description = method
		// If finished, the span will be sent to sentry
		defer span.Finish()

		// AddHook the sentry trace and baggage headers to the metadata
		ctx = getSpanContext(span)

		stream, err := streamer(ctx, desc, cc, method, callOpts...)
		if err != nil && !IsGrpcError(err) {
			hub := sentry.GetHubFromContext(ctx)
			if hub != nil {
				hub.CaptureException(err)
			} else {
				log.Ctx(ctx).Warn().Msg("No sentry hub found in context")
			}
		}

		return stream, err
	}
}

// getSpanContext gets the span context and adds the sentry trace and baggage headers to the metadata
func getSpanContext(span *sentry.Span) context.Context {
	ctx := span.Context()

	// AddHook the sentry trace and baggage headers to the metadata
	trace := span.ToSentryTrace()
	if trace != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, sentry.SentryTraceHeader, trace)
	}
	baggage := span.ToBaggage()
	if baggage != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, sentry.SentryBaggageHeader, baggage)
	}

	return ctx
}
