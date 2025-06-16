package logger_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/paccolamano/goshare/logger"
	"github.com/stretchr/testify/require"
)

func TestHandlerTraceIDInjected(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		format string
	}{
		{"text format", "text"},
		{"json format", "json"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			h := logger.NewTraceHandler(buf, tc.format, "info")
			l := slog.New(h)

			ctx := context.WithValue(t.Context(), logger.TraceIDKey, "abc123")
			l.InfoContext(ctx, "hello world")

			out := buf.String()
			require.NotEmpty(t, out)
			require.Contains(t, out, "traceUUID")
			require.Contains(t, out, "abc123")
		})
	}
}

func TestHandlerDebugSuppressedAtInfo(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	h := logger.NewTraceHandler(buf, "text", "info")
	l := slog.New(h)

	l.Debug("this should not be logged")

	require.Empty(t, buf.String())
}

func TestHandlerDebugEmittedAtDebug(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	h := logger.NewTraceHandler(buf, "text", "debug")
	l := slog.New(h)

	msg := "debug message"
	l.Debug(msg)

	out := buf.String()
	require.NotEmpty(t, out)
	require.Contains(t, out, msg)
}
