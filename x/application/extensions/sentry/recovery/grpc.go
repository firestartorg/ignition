package recovery

import (
	"context"
	"github.com/getsentry/sentry-go"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoverWithContext recovers from the panic `err` by sending it to sentry.
// If the context contains a sentry hub, it will be used, otherwise the current hub will be used.
// If the err is nil, it will call recover() to get the error.
func RecoverWithContext(ctx context.Context, err any) {
	// Get the sentry hub from the context
	var hub *sentry.Hub
	if sentry.HasHubOnContext(ctx) {
		hub = sentry.GetHubFromContext(ctx)
	} else {
		hub = sentry.CurrentHub()
	}

	hub.RecoverWithContext(ctx, err)
}

// RpcRecoveryHandler is a function that recovers from the panic `err` by returning an `error`.
// Compatible with grpc-ecosystem/go-grpc-middleware/recovery package.
func RpcRecoveryHandler(ctx context.Context, err any) error {
	RecoverWithContext(ctx, err)

	return status.Errorf(codes.Internal, "%v", err)
}

var _ grpcrecovery.RecoveryHandlerFuncContext = RpcRecoveryHandler
