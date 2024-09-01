package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jeamon/gobackup/pkg/events"
	"github.com/jeamon/gobackup/pkg/logger"
)

const (
	backupFileExtension = ".bak"
)

// Monitor is an interface defining the behavior of any object
// capable to notifier on the source folder content changes.
type Monitor interface {
	Start(context.Context, <-chan struct{}, events.Queue) error
	Stop() error
}

// App is the structue of an app instance.
type App struct {
	pid       int                  // process id for this App instance.
	srcFolder string               // absolute path of folder to monitor.
	dstFolder string               // backup folder absolute path.
	notifier  Monitor              // concrete object of Monitor contract.
	stop      chan struct{}        // helps goroutines to stop on exit signal.
	jobs      events.Queue         // queue to store instant tasks to handle.
	store     map[string]time.Time // store infos for scheduled deletion action.
	wg        *sync.WaitGroup      // helps ensure all goroutines are stopped.
	mutex     *sync.RWMutex        // mutex to synchronize operations on tasks store.
	log       logger.Logger        // app level json-based logger.
}

// New configures a new App instance.
func New(queueSize int, pid int, src, dst string, monitor Monitor, logger logger.Logger) *App {
	return &App{
		pid:       pid,
		srcFolder: src,
		dstFolder: dst,
		notifier:  monitor,
		stop:      make(chan struct{}, 1),
		jobs:      make(events.Queue, queueSize),
		store:     make(map[string]time.Time),
		wg:        &sync.WaitGroup{},
		mutex:     &sync.RWMutex{},
		log:       logger,
	}
}

// sigHandler handles syscall signals required to stop the program.
func (app *App) sigHandler(sigChan chan os.Signal) {
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT,
		syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)

	<-sigChan
	signal.Stop(sigChan)
	close(app.stop)
}

// monitorFiles calls the monitoring routine of the App instance watcher in
// order to start gathering events of each watched files and errors.
func (app *App) monitorFiles(ctx context.Context) error {
	return app.notifier.Start(ctx, app.stop, app.jobs)
}

// Stop stops the app instance by closing
// the stop channel which all goroutines
// and workers listen on to exit.
func (app *App) Stop() {
	close(app.stop)
}

// CloseQueue close the channel of events.
func (app *App) CloseQueue() {
	close(app.jobs)
}

// start prepares and performs all required routines needed
// to watch and monitor files from source folder.
func (app *App) start(maxWorkers int) (int, error) {
	ctx := context.Background()
	sigChan := make(chan os.Signal, 1)
	go app.sigHandler(sigChan)
	app.startDeleteWorker()
	app.startBackupWorkers(maxWorkers)
	err := app.monitorFiles(ctx)
	if err != nil {
		return 1, fmt.Errorf("failed to start files monitor: %v", err)
	}
	app.CloseQueue()
	app.wg.Wait()
	err = app.SaveAsZipFile(time.Now().UTC())
	if err != nil {
		return 1, err
	}
	return 0, nil
}
