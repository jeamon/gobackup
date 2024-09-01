package app

import (
	"fmt"
	"os"

	"github.com/jeamon/gobackup/pkg/logger"
	"github.com/jeamon/gobackup/pkg/notifier"
	"github.com/jeamon/gobackup/pkg/utils"
	"github.com/jeamon/gobackup/pkg/viewer"
	"github.com/jeamon/gorsn"
)

// Ensure that `*notifier.Notifier` always implements Monitor interface.
var _ Monitor = (*notifier.Notifier)(nil)

// Backup finalizes the initialization of an App instance and
// orchestrates required routines to monitor and handle changes.
func Backup(maxWorkers int, logfile, src, dst, commit, tag string) (int, error) {
	if !utils.IsDirPath(src) || !utils.IsDirPath(dst) {
		return 1, fmt.Errorf("invalid source or backup folder paths. run --help for usage")
	}

	file, logger, err := logger.New(logfile, commit, tag, os.Getpid())
	if err != nil {
		return 1, fmt.Errorf("failed to setup logger: %v", err)
	}
	defer file.Close()
	var opts gorsn.Options
	notifier, err := notifier.New(src, &opts)
	if err != nil {
		return 1, fmt.Errorf("backup: %v", err)
	}
	app := New(maxWorkers, os.Getpid(), src, dst, notifier, logger)
	return app.start(maxWorkers)
}

// ViewLogs uses logview routines to process the content
// of a given log file based on provided filters.
func ViewLogs(logfile, date, reg string) (int, error) {
	if !viewer.IsValidFilters(date, reg) {
		return 1, fmt.Errorf("invalid date and/or regex")
	}
	file, err := viewer.Open(logfile)
	if err != nil {
		return 1, fmt.Errorf("cannot open file: %w", err)
	}
	return viewer.Filter(file, date, reg)
}
