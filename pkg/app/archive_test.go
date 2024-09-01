package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jeamon/gobackup/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetZipID(t *testing.T) {
	app := &App{pid: 1111}
	dt, err := time.Parse(time.RFC3339, "2023-08-14T10:00:00Z")
	require.NoError(t, err)
	id := app.getZipID(dt)
	assert.Equal(t, "20230814.100000.1111", id)
}

func TestSave(t *testing.T) {
	folder, err := os.MkdirTemp("", "folder")
	require.NoError(t, err)
	defer os.RemoveAll(folder)

	dst, err := os.MkdirTemp(folder, "backup")
	require.NoError(t, err)

	file, err := os.CreateTemp(dst, "file")
	require.NoError(t, err)
	file.Close()

	app := &App{dstFolder: dst}
	id := "20230814.100000.1111"
	t.Run("success", func(t *testing.T) {
		success, fails, msg, path, err := app.save(id)
		require.NoError(t, err)
		assert.Equal(t, 1, success)
		assert.Equal(t, 0, fails)
		assert.Equal(t, "success: save backup folder state", msg)
		zipFilename := fmt.Sprintf("%s.%s.zip", filepath.Base(dst), id)
		require.Equal(t, filepath.Join(folder, zipFilename), path)
		assert.FileExists(t, path)
	})

	t.Run("fail", func(t *testing.T) {
		app.dstFolder = filepath.Join(app.dstFolder, "noexist.folderpath")
		success, fails, msg, path, err := app.save(id)
		require.Error(t, err)
		assert.Equal(t, 0, success)
		assert.Equal(t, 0, fails)
		assert.Equal(t, "failed: load backup files", msg)
		require.Equal(t, app.dstFolder, path)
	})
}

func TestSaveAsZipFile(t *testing.T) {
	folder, err := os.MkdirTemp("", "folder")
	require.NoError(t, err)
	defer os.RemoveAll(folder)

	dst, err := os.MkdirTemp(folder, "backup")
	require.NoError(t, err)

	file, err := os.CreateTemp(dst, "file")
	require.NoError(t, err)
	file.Close()

	dt, err := time.Parse(time.RFC3339, "2023-08-14T10:00:00Z")
	require.NoError(t, err)
	out := bytes.NewBuffer(nil)
	app := &App{pid: 1111, dstFolder: dst, log: testhelpers.NewTestLogger(t, out)}
	expectedZipFilename := fmt.Sprintf("%s.%s.zip", filepath.Base(dst), "20230814.100000.1111")
	t.Run("success", func(t *testing.T) {
		app.SaveAsZipFile(dt)
		var data map[string]interface{}
		err := json.Unmarshal(out.Bytes(), &data)
		require.NoError(t, err)

		level, ok := data["level"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "INFO", level)

		msg, ok := data["msg"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "success: save backup folder state [success/fails: 1/0]", msg)

		event, ok := data["event"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "SAVE", event)

		path, ok := data["path"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, filepath.Join(folder, expectedZipFilename), path)
	})

	t.Run("fail", func(t *testing.T) {
		out.Reset()
		app.dstFolder = filepath.Join(app.dstFolder, "noexist.folderpath")
		app.SaveAsZipFile(dt)
		var data map[string]interface{}
		err := json.Unmarshal(out.Bytes(), &data)
		require.NoError(t, err)

		level, ok := data["level"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "ERROR", level)

		msg, ok := data["msg"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "failed: load backup files [success/fails: 0/0]", msg)

		event, ok := data["event"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "SAVE", event)

		path, ok := data["path"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, app.dstFolder, path)
	})
}
