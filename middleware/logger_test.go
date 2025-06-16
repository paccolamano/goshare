package middleware_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paccolamano/goshare/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoggerMiddleware(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l := NewMockInfoLogger(ctrl)
	gomock.InOrder(
		l.EXPECT().
			InfoContext(gomock.Any(),
				"incoming request",
				slog.String("method", http.MethodGet),
				slog.String("path", "/test"),
				slog.String("query", "param=value"),
				slog.String("ip", "192.168.0.1"),
				slog.String("userAgent", "test-agent"),
				slog.Int64("contentLength", 0),
			).
			Times(1),
		l.EXPECT().
			InfoContext(gomock.Any(),
				"request completed",
				slog.String("method", http.MethodGet),
				slog.String("path", "/test"),
				slog.Int("status", http.StatusTeapot),
				gomock.Any(),
			).
			Times(1),
	)

	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true

		w.WriteHeader(http.StatusTeapot)

		if _, err := w.Write([]byte("I'm a teapot")); err != nil {
			t.Log(err)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/test?param=value", nil)

	req.Header.Set("X-Real-IP", "192.168.0.1")
	req.Header.Set("User-Agent", "test-agent")

	w := httptest.NewRecorder()

	middleware.Logger(middleware.WithLogger(l))(handler).ServeHTTP(w, req)

	resp := w.Result()
	body := w.Body.String()

	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
	assert.Equal(t, "I'm a teapot", body)
	assert.True(t, called)
}
