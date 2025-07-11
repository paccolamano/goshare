// Code generated by MockGen. DO NOT EDIT.
// Source: logger.go
//
// Generated by this command:
//
//	mockgen -source=logger.go -destination=logger_mock_test.go -package=middleware_test
//

// Package middleware_test is a generated GoMock package.
package middleware_test

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockInfoLogger is a mock of InfoLogger interface.
type MockInfoLogger struct {
	ctrl     *gomock.Controller
	recorder *MockInfoLoggerMockRecorder
	isgomock struct{}
}

// MockInfoLoggerMockRecorder is the mock recorder for MockInfoLogger.
type MockInfoLoggerMockRecorder struct {
	mock *MockInfoLogger
}

// NewMockInfoLogger creates a new mock instance.
func NewMockInfoLogger(ctrl *gomock.Controller) *MockInfoLogger {
	mock := &MockInfoLogger{ctrl: ctrl}
	mock.recorder = &MockInfoLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInfoLogger) EXPECT() *MockInfoLoggerMockRecorder {
	return m.recorder
}

// InfoContext mocks base method.
func (m *MockInfoLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, msg}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "InfoContext", varargs...)
}

// InfoContext indicates an expected call of InfoContext.
func (mr *MockInfoLoggerMockRecorder) InfoContext(ctx, msg any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, msg}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InfoContext", reflect.TypeOf((*MockInfoLogger)(nil).InfoContext), varargs...)
}
