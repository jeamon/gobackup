package app

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"time"
)

// defines ops name for saving backup folder.
const SAVE string = "SAVE"

// getZipID builds the suffix based on datetime provided and the app process id.
// This value is used to build a unique filename for zip archive of the backup folder.
func (app *App) getZipID(now time.Time) string {
	return fmt.Sprintf("%d%02d%02d.%02d%02d%02d.%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), app.pid)
}

// save creates a zip archive of backup folder. It returns the statistics like
// number of successful files added and the number of failures along with the
// details (message and path and error if any) to insert a log entry.
func (app *App) save(zipID string) (success, fails int, msg, path string, err error) {
	zipFilepath := fmt.Sprintf("%s.%s.zip", app.dstFolder, zipID)
	zfile, err := os.Create(zipFilepath)
	if err != nil {
		msg = "failed: create zip file"
		path = zipFilepath
		return
	}
	defer zfile.Close()

	zw := zip.NewWriter(zfile)
	defer zw.Close()

	files, err := os.ReadDir(app.dstFolder)
	if err != nil {
		msg = "failed: load backup files"
		path = app.dstFolder
		return
	}

	err = os.Chdir(app.dstFolder)
	if err != nil {
		msg = "failed: change working directory to backup folder"
		path = app.dstFolder
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		f, zerr := os.Open(file.Name())
		if zerr != nil {
			fails++
			continue
		}

		w, zerr := zw.Create(file.Name())
		if zerr != nil {
			fails++
			continue
		}

		if _, zerr := io.Copy(w, f); zerr != nil {
			fails++
			continue
		}

		f.Close()
		success++
	}
	msg = "success: save backup folder state"
	path = zipFilepath
	return success, fails, msg, path, err
}

// SaveAsZipFile orchestrates the creation of a zip archive of backup folder.
func (app *App) SaveAsZipFile(t time.Time) error {
	zipID := app.getZipID(t)
	success, fails, msg, path, err := app.save(zipID)
	if err != nil {
		app.log.Error(fmt.Sprintf("%s [success/fails: %d/%d]", msg, success, fails), SAVE, path, err)
		return err
	}
	app.log.Info(fmt.Sprintf("%s [success/fails: %d/%d]", msg, success, fails), SAVE, path)
	return nil
}
