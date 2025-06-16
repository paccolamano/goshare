package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paccolamano/goshare/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRecoverWithPanic(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l := NewMockErrorLogger(ctrl)
	l.EXPECT().
		ErrorContext(gomock.Any(),
			"recovered from panic",
			gomock.Any(),
		).
		Times(1)

	handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("something went wrong")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	middleware.Recover(middleware.WithRecoverLogger(l))(handler).ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	expectedBody := `{"error":"Internal Server Error"}` + "\n"
	body := w.Body.String()
	assert.Equal(t, expectedBody, body)
}

func TestRecoverNoPanic(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l := NewMockErrorLogger(ctrl)
	l.EXPECT().
		ErrorContext(gomock.Any(),
			"recovered from panic",
			gomock.Any(),
		).
		Times(0)

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		handlerCalled = true

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	w := httptest.NewRecorder()

	middleware.Recover(middleware.WithRecoverLogger(l))(handler).ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, handlerCalled)
}
