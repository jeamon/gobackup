package app

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jeamon/gobackup/pkg/events"
	"github.com/jeamon/gobackup/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartBackupWorkers(t *testing.T) {
	max := 5
	app := New(1, 0, "", "", nil, nil)
	app.startBackupWorkers(max)
	done := false
	go func() {
		app.wg.Wait()
		done = true
	}()
	for i := 0; i < max; i++ {
		app.wg.Done()
	}
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, true, done)
}

func TestStartDeleteWorker(t *testing.T) {
	folder, err := os.MkdirTemp("", "folder")
	require.NoError(t, err)
	defer os.RemoveAll(folder)

	t.Run("success deletion", func(t *testing.T) {
		file, err := os.CreateTemp(folder, "file")
		require.NoError(t, err)
		filePath := file.Name()
		file.Close()

		app := New(0, 0, "", "", nil, testhelpers.NewTestLogger(t, io.Discard))
		app.store[filePath] = time.Now().Add(-500 * time.Millisecond)
		app.startDeleteWorker()
		time.Sleep(600 * time.Millisecond)
		assert.NoFileExists(t, filePath)
		app.mutex.RLock()
		num := len(app.store)
		app.mutex.RUnlock()
		assert.Equal(t, 0, num)
		close(app.stop)
	})
	t.Run("deletion time not reached", func(t *testing.T) {
		file, err := os.CreateTemp(folder, "file")
		require.NoError(t, err)
		filePath := file.Name()
		file.Close()
		app := New(0, 0, "", "", nil, testhelpers.NewTestLogger(t, io.Discard))
		app.store[filePath] = time.Now().Add(1 * time.Minute)
		app.startDeleteWorker()
		// making sure at least first deletion round did run.
		time.Sleep(500 * time.Millisecond)
		assert.FileExists(t, filePath)
		app.mutex.RLock()
		num := len(app.store)
		app.mutex.RUnlock()
		assert.Equal(t, 1, num)
		close(app.stop)
	})

	t.Run("fail to delete", func(t *testing.T) {
		app := New(0, 0, "", "", nil, testhelpers.NewTestLogger(t, io.Discard))
		app.store["filepath.noexist"] = time.Now().Add(-500 * time.Millisecond)
		app.startDeleteWorker()
		// making sure at least first deletion round did run.
		time.Sleep(600 * time.Millisecond)
		app.mutex.RLock()
		num := len(app.store)
		app.mutex.RUnlock()
		assert.Equal(t, 0, num)
		close(app.stop)
	})
}

func TestBackupWorker(t *testing.T) {
	folder, err := os.MkdirTemp("", "folder")
	require.NoError(t, err)
	defer os.RemoveAll(folder)

	src, err := os.MkdirTemp(folder, "source")
	require.NoError(t, err)
	dst, err := os.MkdirTemp(folder, "backup")
	require.NoError(t, err)

	t.Run("event:create", func(t *testing.T) {
		sfile, err := os.CreateTemp(src, "sfile")
		require.NoError(t, err)
		sfilePath := sfile.Name()
		sfile.Close()
		app := New(1, 0, src, dst, nil, testhelpers.NewTestLogger(t, io.Discard))
		go func() {
			app.jobs <- &events.Change{Path: sfilePath, Ops: events.CREATE}
			close(app.stop)
		}()
		app.wg.Add(1)
		app.backupWorker(1)
		assert.FileExists(t, filepath.Join(dst, filepath.Base(sfilePath)+backupFileExtension))
	})
}
