package logger

import (
	"context"
	"io"
	"log/slog"
)

// contextKey is a custom type to avoid context key collisions.
type contextKey string

// TraceIDKey is the context key used to store and retrieve the trace ID.
// If present, the trace ID will be included as a log attribute.
const TraceIDKey contextKey = "traceUUID"

// TraceHandler wraps slog.Handler and injects the trace ID from context into log records.
type TraceHandler struct {
	slog.Handler
}

// NewTraceHandler creates a new Handler with the given output writer, format, and log level.
//
// Parameters:
//   - out: the io.Writer where logs will be written.
//   - format: output format, either "text" or "json" (default is "text").
//   - level: log level string, can be "debug", "warn", "error" or any other value (default is "info").
//
// Returns:
//   - A Handler that wraps the appropriate slog.Handler with the specified configuration.
//
// If the format is "json", the handler will output logs in JSON format.
// The log level controls the minimum level of logs emitted.
// If level is "debug", the handler will also include source file information.
func NewTraceHandler(out io.Writer, format string, level string) *TraceHandler {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}

	switch level {
	case "debug":
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	}

	handler := &TraceHandler{Handler: slog.NewTextHandler(out, opts)}
	if format == "json" {
		handler.Handler = slog.NewJSONHandler(out, opts)
	}

	return handler
}

// Handle adds the trace ID from the context to the log record (if available)
// and delegates the log handling to the wrapped slog.Handler.
//
// Parameters:
//   - ctx: context potentially containing a trace ID under TraceIDKey.
//   - r: the slog.Record to be handled.
//
// Returns:
//   - An error if the underlying handler returns an error.
func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	if v, ok := ctx.Value(TraceIDKey).(string); ok {
		r.AddAttrs(slog.String("traceUUID", v))
	}

	return h.Handler.Handle(ctx, r)
}

// WithAttrs returns a new TraceHandler whose attributes consists
// of h's attributes followed by attrs.
func (h *TraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TraceHandler{Handler: h.Handler.WithAttrs(attrs)}
}
