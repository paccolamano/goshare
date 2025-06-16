package middleware

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"
)

//go:generate mockgen -source=logger.go -destination=logger_mock_test.go -package=middleware_test

// InfoLogger defines the interface required for structured logging.
// It is compatible with log/slog or any other logger implementing
// InfoContext for contextual logging.
type InfoLogger interface {
	// InfoContext logs an informational message with context and key-value pairs.
	InfoContext(ctx context.Context, msg string, args ...any)
}

// RecoverOptions holds configuration options for the Recover middleware,
// including the logger and the panic callback handler.
type LoggerOptions struct {
	Logger InfoLogger
}

// RecoverOption represents a functional option for configuring Recover middleware.
type LoggerOption func(*LoggerOptions)

// WithLogger sets a custom InfoLogger for the Logger middleware.
// It allows integration with custom or structured loggers.
func WithLogger(l InfoLogger) LoggerOption {
	return func(opt *LoggerOptions) {
		opt.Logger = l
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger returns an HTTP middleware that logs each incoming request
// and its corresponding response. It logs request method, path, query,
// client IP, user agent, content length, response status, and duration.
//
// Example usage:
//
//	http.Handle("/path", Logger(logger)(yourHandler))
//
// The logger must implement the InfoLogger interface, typically wrapping
// a structured logger such as slog.Logger.
func Logger(opts ...LoggerOption) func(http.Handler) http.Handler {
	options := &LoggerOptions{
		Logger: slog.Default(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			ip := r.Header.Get("X-Real-IP")
			if ip == "" {
				ip, _, _ = net.SplitHostPort(r.RemoteAddr)
			}

			options.Logger.InfoContext(r.Context(), "incoming request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.String("ip", ip),
				slog.String("userAgent", r.UserAgent()),
				slog.Int64("contentLength", r.ContentLength),
			)

			next.ServeHTTP(rw, r)

			options.Logger.InfoContext(r.Context(), "request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rw.statusCode),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}
