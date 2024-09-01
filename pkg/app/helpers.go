package app

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jeamon/gobackup/pkg/utils"
)

// UpdateBackupFileContent copies the content of a given file path
// to its the backup file.
func (app *App) UpdateBackupFileContent(path string) (err error) {
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.OpenFile(filepath.Join(app.dstFolder, filepath.Base(path))+backupFileExtension, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	defer func() {
		if c := w.Close(); c != nil && err == nil {
			err = c
		}
	}()

	_, err = io.Copy(w, r)
	return err
}

// CreateBackupFile creates a file into the backup folder with
// same name as the original file and use `.bat` as extension.
func (app *App) CreateBackupFile(path string) error {
	f, err := os.Create(filepath.Join(app.dstFolder, filepath.Base(path)) + backupFileExtension)
	if err != nil {
		return err
	}
	return f.Close()
}

// IsImmediateDelete checks wether the filename matches the required pattern
// to trigger immediate deletion action of that file.
func (app *App) IsImmediateDelete(path string) bool {
	return strings.HasPrefix(path, "delete_") && len(strings.TrimPrefix(path, "delete_")) > 0
}

// IsScheduleDelete checks if a filename matches the required pattern
// to trigger a scheduling deletion action of that file. It matches if
// the name follow this pattern `delete_ISODATETIME`. The ISODATETIME
// should matches the RFC3339 ("2006-01-02T15:04:05Z07:00") format.
// If so, it returns the absolute source & backup filepaths and datetime.
func (app *App) IsScheduleDelete(path string) (string, string, time.Time, bool) {
	if !strings.HasPrefix(filepath.Base(path), "delete_") {
		return "", "", time.Time{}, false
	}

	suffix := strings.TrimPrefix(filepath.Base(path), "delete_")
	isodatetime, filename, found := strings.Cut(suffix, "_")
	if found && len(filename) > 0 {
		if at, err := time.Parse(time.RFC3339, utils.FixColonCharacter(isodatetime)); err == nil {
			spath := filepath.Join(filepath.Dir(path), filename)
			dpath := filepath.Join(app.dstFolder, filename) + backupFileExtension
			return spath, dpath, at, true
		}
	}

	return "", "", time.Time{}, false
}
