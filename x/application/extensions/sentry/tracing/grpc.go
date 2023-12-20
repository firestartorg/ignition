package tracing

import (
	"context"
	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ContinueFromRpcRequest(ctx context.Context) sentry.SpanOption {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return func(s *sentry.Span) {}
	}

	return sentry.ContinueFromHeaders(getHeader(md, sentry.SentryTraceHeader), getHeader(md, sentry.SentryBaggageHeader))
}

func getHeader(md metadata.MD, key string) string {
	values, ok := md[key]
	if !ok {
		return ""
	}
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func ToSpanStatus(code codes.Code) sentry.SpanStatus {
	switch code {
	case codes.OK:
		return sentry.SpanStatusOK
	case codes.Canceled:
		return sentry.SpanStatusCanceled
	case codes.Unknown:
		return sentry.SpanStatusUnknown
	case codes.InvalidArgument:
		return sentry.SpanStatusInvalidArgument
	case codes.DeadlineExceeded:
		return sentry.SpanStatusDeadlineExceeded
	case codes.NotFound:
		return sentry.SpanStatusNotFound
	case codes.AlreadyExists:
		return sentry.SpanStatusAlreadyExists
	case codes.PermissionDenied:
		return sentry.SpanStatusPermissionDenied
	case codes.ResourceExhausted:
		return sentry.SpanStatusResourceExhausted
	case codes.FailedPrecondition:
		return sentry.SpanStatusFailedPrecondition
	case codes.Aborted:
		return sentry.SpanStatusAborted
	case codes.OutOfRange:
		return sentry.SpanStatusOutOfRange
	case codes.Unimplemented:
		return sentry.SpanStatusUnimplemented
	case codes.Internal:
		return sentry.SpanStatusInternalError
	case codes.Unavailable:
		return sentry.SpanStatusUnavailable
	case codes.DataLoss:
		return sentry.SpanStatusDataLoss
	case codes.Unauthenticated:
		return sentry.SpanStatusUnauthenticated
	default:
		return sentry.SpanStatusUndefined
	}
}

func IsGrpcError(err error) bool {
	_, ok := status.FromError(err)
	return ok
}
