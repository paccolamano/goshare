package syncute_test

import (
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/paccolamano/goshare/syncute"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/mock/gomock"
)

func TestWait(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup

	called := false

	syncute.Wait(&wg, func() {
		called = true
	})

	wg.Wait()

	assert.True(t, called)
}

func TestRunWithShutdown(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := NewMockDebugLogger(ctrl)
	service1 := NewMockService(ctrl)
	service2 := NewMockService(ctrl)

	service1.EXPECT().Run(gomock.Any()).Times(1)
	service2.EXPECT().Run(gomock.Any()).Times(1)

	service1.EXPECT().Shutdown(gomock.Any()).Times(1)
	service2.EXPECT().Shutdown(gomock.Any()).Times(1)

	logger.EXPECT().DebugContext(gomock.Any(), "shutdown signal received", "timeout", "1s").Times(1)
	logger.EXPECT().DebugContext(gomock.Any(), mock.MatchedBy(func(msg string) bool {
		return msg == "graceful shutdown completed" || msg == "forced shutdown: timeout reached"
	})).Times(1)

	// Run the service runner with a short-lived context
	go func(tb testing.TB) {
		tb.Helper()
		// Allow RunWithShutdown to start
		time.Sleep(50 * time.Millisecond)
		// Send termination signal manually
		if err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM); err != nil {
			tb.Errorf("failed to send sigterm to process %d", syscall.Getpid())
		}
	}(t)

	syncute.RunWithShutdown(logger, 1*time.Second, service1, service2)
}
