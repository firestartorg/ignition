package tracing

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func HttpMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// AddHook a sentry hub to the request context
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			log.Ctx(ctx).Warn().Msg("No sentry hub found in context")
			handler.ServeHTTP(w, r)
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

		// AddHook the transaction to the request context
		r = r.WithContext(transaction.Context())
		hub.Scope().SetRequest(r) // AddHook the request to the scope

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
