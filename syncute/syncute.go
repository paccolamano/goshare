package syncute

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

//go:generate mockgen -source=syncute.go -destination=syncute_mock_test.go -package=syncute_test

// Wait wraps a function execution within a goroutine and increments the given sync.WaitGroup.
//
// It ensures that the function `f` is executed in a separate goroutine,
// and that `wg.Done()` is called when `f` completes.
//
// Parameters:
//   - wg: pointer to a sync.WaitGroup to track the goroutine.
//   - f: function to execute concurrently.
func Wait(wg *sync.WaitGroup, f func()) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		f()
	}()
}

// DebugLogger is an interface for logging debug messages with contextual information.
type DebugLogger interface {
	// DebugContext logs a debug message with the provided context and optional arguments.
	DebugContext(ctx context.Context, msg string, args ...any)
}

// Service represents a long-running service with support for graceful shutdown.
type Service interface {
	// Run starts the service and should block until the service is stopped or the context is cancelled.
	Run(ctx context.Context)

	// Shutdown is called when a shutdown signal is received. It should clean up resources and terminate gracefully.
	Shutdown(ctx context.Context)
}

// RunWithShutdown executes multiple services concurrently and gracefully shuts them down upon receiving
// a termination signal (SIGINT or SIGTERM).
//
// It listens for OS signals to initiate shutdown, runs all services concurrently, and waits for all of them to complete.
// If the shutdown phase exceeds the given timeout, it logs a forced shutdown.
//
// Parameters:
//   - logger: a DebugLogger implementation for logging lifecycle events.
//   - timeout: the duration allowed for graceful shutdown before a forced termination.
//   - services: a variadic list of Service implementations to manage.
func RunWithShutdown(logger DebugLogger, timeout time.Duration, services ...Service) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	for _, v := range services {
		Wait(&wg, func() {
			v.Run(ctx)
		})
	}

	<-ctx.Done()
	logger.DebugContext(ctx, "shutdown signal received", "timeout", timeout.String())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, v := range services {
		Wait(&wg, func() {
			v.Shutdown(shutdownCtx)
		})
	}

	done := make(chan bool, 1)
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.DebugContext(context.Background(), "graceful shutdown completed")
	case <-shutdownCtx.Done():
		logger.DebugContext(context.Background(), "forced shutdown: timeout reached")
	}
}
