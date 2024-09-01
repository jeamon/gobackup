package testhelpers

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/jeamon/gobackup/pkg/events"
	"github.com/jeamon/gobackup/pkg/logger"
)

// MockMonitor represents a mock of Monitor interface.
type MockMonitor struct {
	StopFunc  func() error
	StartFunc func(context.Context, <-chan struct{}, events.Queue) error
}

// AddFolder mocks the behavior of folder adding by the watcher.
func (m *MockMonitor) Stop() error {
	return m.StopFunc()
}

// Start mocks the behavior of watcher running.
func (m *MockMonitor) Start(ctx context.Context, quit <-chan struct{}, jobs events.Queue) error {
	return m.StartFunc(ctx, quit, jobs)
}

// NewTestLogger provides a logger for testing purposes. For Noop use `io.Discard`
// as output. To hold the log entries, you can use a pointer of `bytes.Buffer`.
func NewTestLogger(t *testing.T, out io.Writer) logger.Logger {
	t.Helper()
	return &logger.DefaultLogger{Log: slog.New(slog.NewJSONHandler(out, nil)).With(slog.String("commit", ""), slog.String("tag", ""), slog.Int("pid", 0))}
}
