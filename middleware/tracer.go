package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// TracerOptions holds configuration options for the Tracer middleware.
type TracerOptions struct {
	TraceKey any
}

// TracerOption represents a functional option for configuring Tracer middleware.
type TracerOption func(*TracerOptions)

// WithTraceKey sets a custom trace key for the Tracer middleware.
func WithTraceKey(key any) TracerOption {
	return func(opt *TracerOptions) {
		opt.TraceKey = key
	}
}

// Tracer returns a middleware that generates a unique request ID (UUID) for each incoming HTTP request,
// attaches it to the response header as "X-Request-ID", and stores it in the request context using the provided key.
//
// This is useful for request tracing, correlation across distributed systems, and contextual logging.
//
// Parameters:
//   - traceIDKey: the context key under which the generated UUID will be stored in the request context.
//
// Example usage:
//
//	http.Handle("/api", Tracer("traceID")(yourHandler))
//
// You can then retrieve the trace ID later in the request lifecycle:
//
//	traceID := r.Context().Value("traceID").(string)
func Tracer(opts ...TracerOption) func(http.Handler) http.Handler {
	options := &TracerOptions{
		TraceKey: "traceUUID",
	}

	for _, opt := range opts {
		opt(options)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uuid := uuid.New().String()
			w.Header().Set("X-Request-ID", uuid)
			ctx := context.WithValue(r.Context(), options.TraceKey, uuid)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
