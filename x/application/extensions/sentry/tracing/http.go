package tracing

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"net/http"
	"time"
)

func NewHttpMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Add a sentry hub to the request context
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			return
		}

		options := []sentry.SpanOption{
			sentry.WithOpName("http.server"),
			sentry.ContinueFromRequest(r),
			sentry.WithTransactionSource(sentry.SourceURL),
		}

		// Start a new transaction
		transaction := sentry.StartTransaction(ctx,
			fmt.Sprintf("%s %s", r.Method, r.URL.Path),
			options...,
		)
		defer transaction.Finish()

		// Add the transaction to the request context
		r = r.WithContext(transaction.Context())
		hub.Scope().SetRequest(r) // Add the request to the scope

		// Recover from panics and send them to sentry
		defer func() {
			if err := recover(); err != nil {
				eventID := hub.RecoverWithContext(
					context.WithValue(r.Context(), sentry.RequestContextKey, r),
					err,
				)
				if eventID != nil {
					hub.Flush(2 * time.Second)
				}
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
