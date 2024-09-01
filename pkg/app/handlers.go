package app

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jeamon/gobackup/pkg/events"
	"github.com/jeamon/gobackup/pkg/utils"
)

// ReceiveFolderEventHandler handles folder added into the source
// folder. It configures the App watcher to immediately monitor all
// files inside the folder and its sub-folders. All files found are
// backed up immediately according to events rules.
func (app *App) ReceiveFolderEventHandler(path string) {
	err := utils.CreateFolder(path)
	if err != nil {
		app.log.Error("failed: create folder", string(events.WATCH), path, err)
		return
	}
	app.log.Info("success: create folder", string(events.WATCH), path)
}

// DeleteRequestHandler orchestrates the processing of file deletion request.
func (app *App) DeleteRequestHandler(path string) {
	var err error
	filename := strings.TrimPrefix(filepath.Base(path), "delete_")
	spath := filepath.Join(filepath.Dir(path), filename)
	if err = utils.DeleteFile(spath); err != nil {
		app.log.Error("failed: delete file", string(events.RDELETE), spath, err)
	} else {
		app.log.Info("success: delete file", string(events.RDELETE), spath)
	}

	dpath := filepath.Join(app.dstFolder, filename) + backupFileExtension
	if err = utils.DeleteFile(dpath); err != nil {
		app.log.Error("failed: delete file", string(events.RDELETE), dpath, err)
	} else {
		app.log.Info("success: delete file", string(events.RDELETE), dpath)
	}

	if err = utils.DeleteFile(path); err != nil {
		app.log.Error("failed: delete file", string(events.RDELETE), path, err)
	} else {
		app.log.Info("success: delete file", string(events.RDELETE), path)
	}
}

// CreateEventHandler orchestrates the processing of file creation events.
func (app *App) CreateEventHandler(path string) {
	err := app.CreateBackupFile(path)
	if err != nil {
		app.log.Error("failed: create file", string(events.CREATE), path, err)
		return
	}
	app.log.Info("success: create file", string(events.CREATE), path)
}

// ModifyEventHandler orchestrates the processing of file content modification events.
func (app *App) ModifyEventHandler(path string) {
	err := app.UpdateBackupFileContent(path)
	if err != nil {
		app.log.Error("failed: update file", string(events.MODIFY), path, err)
		return
	}
	app.log.Info("success: update file", string(events.MODIFY), path)
}

// ScheduleDeleteRequests adds each absolute filepath and its deletion
// datetime to the map store.
func (app *App) ScheduleDeleteRequests(at time.Time, paths ...string) {
	app.mutex.Lock()
	for _, path := range paths {
		app.store[path] = at
	}
	app.mutex.Unlock()
}

// RenameEventHandler just logs folder or file renaming events.
func (app *App) RenameEventHandler(path string) {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		app.log.Info("receive: rename folder event", string(events.RENAME), path)
		return
	}
	app.log.Info("receive: rename file event", string(events.RENAME), path)
}

// DeleteEventHandler just logs folder or file delete events.
func (app *App) DeleteEventHandler(path string) {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		app.log.Info("receive: delete folder event", string(events.DELETE), path)
		return
	}
	app.log.Info("receive: delete file event", string(events.DELETE), path)
}

// AttributeEventHandler just logs folder or file attributes events.
func (app *App) AttributeEventHandler(path string) {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		app.log.Info("receive: folder attribute event", string(events.ATTRIBUTE), path)
		return
	}
	app.log.Info("receive: file attribute event", string(events.ATTRIBUTE), path)
}
