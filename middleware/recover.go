package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

//go:generate mockgen -source=recover.go -destination=recover_mock_test.go -package=middleware_test

// ErrorLogger defines an interface for logging errors with contextual information.
// This is typically implemented by structured logging libraries.
type ErrorLogger interface {
	ErrorContext(ctx context.Context, msg string, args ...any)
}

// RecoverOptions holds configuration options for the Recover middleware,
// including the logger and the panic callback handler.
type RecoverOptions struct {
	Logger   ErrorLogger
	Callback func(w http.ResponseWriter, r *http.Request)
}

// RecoverOption represents a functional option for configuring Recover middleware.
type RecoverOption func(*RecoverOptions)

// WithRecoverLogger sets a custom ErrorLogger for the Recover middleware.
// It allows integration with custom or structured loggers.
func WithRecoverLogger(l ErrorLogger) RecoverOption {
	return func(ro *RecoverOptions) {
		ro.Logger = l
	}
}

// WithCallback sets a custom callback function that will be executed
// when a panic is recovered in an HTTP handler.
func WithCallback(f func(w http.ResponseWriter, r *http.Request)) RecoverOption {
	return func(ro *RecoverOptions) {
		ro.Callback = f
	}
}

// Recover returns a middleware that recovers from panics during HTTP request handling.
// It logs the panic using the provided logger (or a default one) and then calls the
// specified callback (or a default 500 error response).
//
// Example usage:
//
//	http.Handle("/api", Recover(WithLogger(myLogger), WithCallback(myCallback))(myHandler))
//
// If no options are provided, it uses slog.Default() as the logger and writes a generic
// 500 Internal Server Error response in JSON format.
func Recover(opts ...RecoverOption) func(http.Handler) http.Handler {
	options := &RecoverOptions{
		Logger: slog.Default(),
		Callback: func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": http.StatusText(http.StatusInternalServerError),
			})
		},
	}

	for _, opt := range opts {
		opt(options)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func(ctx context.Context) {
				if err := recover(); err != nil {
					options.Logger.ErrorContext(ctx, "recovered from panic", "error", err, "stack", string(debug.Stack()))
					options.Callback(w, r)
				}
			}(r.Context())

			next.ServeHTTP(w, r)
		})
	}
}
