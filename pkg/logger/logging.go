package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jeamon/gobackup/pkg/utils"
)

type Logger interface {
	Error(msg, event, path string, err error)
	Info(msg, event, path string)
}

type DefaultLogger struct {
	Log *slog.Logger
}

// Error inserts error level log entry. It cleans the
// provided path by fixing the colon character if any.
func (dl *DefaultLogger) Error(msg, event, path string, err error) {
	dl.Log.Error(msg,
		slog.String("error", err.Error()),
		slog.String("event", event),
		slog.String("path", utils.FixColonCharacter(path)),
	)
}

// Info inserts info level log entry. It cleans the
// provided path by fixing the colon character if any.
func (dl *DefaultLogger) Info(msg, event, path string) {
	dl.Log.Info(msg,
		slog.String("event", event),
		slog.String("path", utils.FixColonCharacter(path)),
	)
}

// setupLogger creates or opens the app log file (default to`file.log`) and initialize
// an instance of slog with some predefined attributes for app logging.
func New(filename, commit, tag string, pid int) (*os.File, Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create log file: %v", err)
	}
	logger := slog.New(slog.NewJSONHandler(file, nil)).With(slog.String("commit", commit), slog.String("tag", tag), slog.Int("pid", pid))
	return file, &DefaultLogger{logger}, nil
}
