package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/paccolamano/goshare/middleware"
	"github.com/stretchr/testify/require"
)

func TestTracerMiddleware(t *testing.T) {
	t.Parallel()

	var traceIDInContext string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceIDVal := r.Context().Value("traceUUID")
		traceIDInContext, _ = traceIDVal.(string)

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/tracer", nil)
	w := httptest.NewRecorder()

	middleware.Tracer(middleware.WithTraceKey("traceUUID"))(handler).ServeHTTP(w, req)

	resp := w.Result()
	requestID := resp.Header.Get("X-Request-ID")

	require.NotEmpty(t, requestID, "X-Request-ID header should be set")
	_, err := uuid.Parse(requestID)
	require.NoError(t, err, "X-Request-ID should be a valid UUID")

	require.Equal(t, requestID, traceIDInContext, "Trace ID in context should match X-Request-ID header")
}
