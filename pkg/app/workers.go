package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jeamon/gobackup/pkg/events"
	"github.com/jeamon/gobackup/pkg/utils"
)

// backupWorker processes each event that comes in the `jobs` queue.
// It only handles directory or regular file associated to an event.
// WATCH and RDELETE events do not trigger any actions. Those are added
// to avoid linters warnings.
func (app *App) backupWorker(id int) {
	defer app.wg.Done()
	for {
		select {
		case ce := <-app.jobs:
			switch ce.Ops {
			case events.CREATE:
				fi, err := os.Stat(ce.Path)
				if err == nil && fi.IsDir() {
					app.ReceiveFolderEventHandler(ce.Path)
					continue
				}
				if !fi.Mode().IsRegular() {
					continue
				}
				if !strings.HasPrefix(filepath.Base(ce.Path), "delete_") {
					app.CreateEventHandler(ce.Path)
					continue
				}
				if spath, dpath, at, ok := app.IsScheduleDelete(ce.Path); ok {
					app.ScheduleDeleteRequests(at, spath, dpath, ce.Path)
					continue
				}
				if app.IsImmediateDelete(filepath.Base(ce.Path)) {
					app.DeleteRequestHandler(ce.Path)
				}

			case events.MODIFY:
				if fi, err := os.Stat(ce.Path); err != nil || fi.IsDir() || !fi.Mode().IsRegular() {
					continue
				}
				if app.IsImmediateDelete(filepath.Base(ce.Path)) {
					continue
				}
				app.ModifyEventHandler(ce.Path)

			case events.RENAME:
				app.RenameEventHandler(ce.Path)

			case events.DELETE:
				app.DeleteEventHandler(ce.Path)

			case events.ATTRIBUTE:
				app.AttributeEventHandler(ce.Path)

			case events.WATCH, events.RDELETE:
			}
		case <-app.stop:
			log.Println("stopped backup worker:", id)
			return
		}
	}
}

// startBackupWorkers pre-boots `maxWorkers` number of workers in charge
// of consuming tasks queued on `jobs` and process them.
func (app *App) startBackupWorkers(maxWorkers int) {
	for i := 0; i < maxWorkers; i++ {
		id := i
		app.wg.Add(1)
		go app.backupWorker(id)
	}
}

// startDeleteWorker starts a goroutine which checks the App store/map
// and deletes files which were scheduled to be removed after a given
// datetime. The chech happens every 500 ms.
func (app *App) startDeleteWorker() {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		for {
			select {
			case <-time.After(500 * time.Millisecond):
				app.mutex.Lock()
				for path, t := range app.store {
					if time.Now().Before(t) {
						continue
					}
					err := utils.DeleteFile(path)
					if err != nil {
						app.log.Error("failed: delete file", string(events.DELETE), path, err)
					} else {
						app.log.Info("success: delete file", string(events.DELETE), path)
					}
					delete(app.store, path)
				}
				app.mutex.Unlock()
			case <-app.stop:
				log.Println("stopped delete worker")
				return
			}
		}
	}()
}
