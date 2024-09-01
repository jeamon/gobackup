package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jeamon/gobackup/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteRequestHandler(t *testing.T) {
	src, err := os.MkdirTemp("", "source")
	require.NoError(t, err)
	defer os.RemoveAll(src)

	sfile, err := os.Create(filepath.Join(src, "file"))
	require.NoError(t, err)
	spath := sfile.Name()
	sfile.Close()

	file, err := os.Create(filepath.Join(src, "delete_file"))
	require.NoError(t, err)
	path := file.Name()
	file.Close()

	dst, err := os.MkdirTemp("", "backup")
	require.NoError(t, err)
	defer os.RemoveAll(dst)

	dfile, err := os.Create(filepath.Join(dst, "file.bak"))
	require.NoError(t, err)
	dpath := dfile.Name()
	dfile.Close()

	out := bytes.NewBuffer(nil)
	logger := testhelpers.NewTestLogger(t, out)
	app := New(1, 0, src, dst, nil, logger)
	app.DeleteRequestHandler(path)

	assert.NoFileExists(t, spath)
	assert.NoFileExists(t, dpath)
	assert.NoFileExists(t, path)
	assert.Equal(t, 4, len(strings.Split(out.String(), "\n")))
}

func TestCreateEventHandler(t *testing.T) {
	src, err := os.MkdirTemp("", "source")
	require.NoError(t, err)
	defer os.RemoveAll(src)

	sfile, err := os.Create(filepath.Join(src, "file"))
	require.NoError(t, err)
	spath := sfile.Name()
	sfile.Close()

	dst, err := os.MkdirTemp("", "backup")
	require.NoError(t, err)
	defer os.RemoveAll(dst)

	out := bytes.NewBuffer(nil)
	logger := testhelpers.NewTestLogger(t, out)
	app := New(1, 0, src, dst, nil, logger)
	t.Run("success", func(t *testing.T) {
		app.CreateEventHandler(spath)
		require.FileExists(t, filepath.Join(dst, "file.bak"))

		var data map[string]interface{}
		err = json.Unmarshal(out.Bytes(), &data)
		require.NoError(t, err)

		level, ok := data["level"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "INFO", level)

		msg, ok := data["msg"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "success: create file", msg)

		event, ok := data["event"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "CREATE", event)

		path, ok := data["path"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, spath, path)
	})

	t.Run("fail", func(t *testing.T) {
		out.Reset()
		// delete the backup folder to trigger the backup
		// file creation failure.
		os.RemoveAll(dst)
		app.CreateEventHandler(spath)
		require.NoFileExists(t, filepath.Join(dst, "file.bak"))

		var data map[string]interface{}
		err = json.Unmarshal(out.Bytes(), &data)
		require.NoError(t, err)

		level, ok := data["level"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "ERROR", level)

		msg, ok := data["msg"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "failed: create file", msg)

		event, ok := data["event"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "CREATE", event)

		path, ok := data["path"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, spath, path)

		_, ok = data["error"].(string)
		assert.Equal(t, true, ok)
	})
}

func TestModifyEventHandler(t *testing.T) {
	src, err := os.MkdirTemp("", "source")
	require.NoError(t, err)
	defer os.RemoveAll(src)

	sfile, err := os.Create(filepath.Join(src, "file"))
	require.NoError(t, err)
	spath := sfile.Name()
	_, err = sfile.WriteString("source file content.")
	require.NoError(t, err)
	sfile.Close()

	dst, err := os.MkdirTemp("", "backup")
	require.NoError(t, err)
	defer os.RemoveAll(dst)

	out := bytes.NewBuffer(nil)
	logger := testhelpers.NewTestLogger(t, out)
	app := New(1, 0, src, dst, nil, logger)
	t.Run("success", func(t *testing.T) {
		app.ModifyEventHandler(spath)
		content, err := os.ReadFile(filepath.Join(dst, "file.bak"))
		require.NoError(t, err)
		assert.Equal(t, []byte("source file content."), content)

		var data map[string]interface{}
		err = json.Unmarshal(out.Bytes(), &data)
		require.NoError(t, err)

		level, ok := data["level"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "INFO", level)

		msg, ok := data["msg"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "success: update file", msg)

		event, ok := data["event"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "MODIFY", event)

		path, ok := data["path"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, spath, path)
	})

	t.Run("fail", func(t *testing.T) {
		out.Reset()
		os.Remove(filepath.Join(dst, "file.bak"))
		os.Remove(spath)
		app.ModifyEventHandler(spath)
		require.NoFileExists(t, filepath.Join(dst, "file.bak"))

		var data map[string]interface{}
		err = json.Unmarshal(out.Bytes(), &data)
		require.NoError(t, err)

		level, ok := data["level"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "ERROR", level)

		msg, ok := data["msg"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "failed: update file", msg)

		event, ok := data["event"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "MODIFY", event)

		path, ok := data["path"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, spath, path)

		_, ok = data["error"].(string)
		assert.Equal(t, true, ok)
	})
}

func TestScheduleDeleteRequests(t *testing.T) {
	app := New(1, 0, "", "", nil, nil)
	now := time.Now()
	app.ScheduleDeleteRequests(now, "file/path", "file.bak/path", "delete_file/path")
	require.Equal(t, 3, len(app.store))
	at, ok := app.store["file/path"]
	require.Equal(t, true, ok)
	assert.Equal(t, now, at)
	at, ok = app.store["file.bak/path"]
	require.Equal(t, true, ok)
	assert.Equal(t, now, at)
	at, ok = app.store["delete_file/path"]
	require.Equal(t, true, ok)
	assert.Equal(t, now, at)
}

func TestRenameEventHandler(t *testing.T) {
	out := bytes.NewBuffer(nil)
	logger := testhelpers.NewTestLogger(t, out)
	app := New(1, 0, "", "", nil, logger)
	t.Run("file", func(t *testing.T) {
		file, err := os.CreateTemp("", "file")
		require.NoError(t, err)
		fpath := file.Name()
		file.Close()
		defer os.Remove(file.Name())
		app.RenameEventHandler(fpath)

		var data map[string]interface{}
		err = json.Unmarshal(out.Bytes(), &data)
		require.NoError(t, err)

		level, ok := data["level"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "INFO", level)

		msg, ok := data["msg"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "receive: rename file event", msg)

		event, ok := data["event"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "RENAME", event)

		path, ok := data["path"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, fpath, path)
	})

	t.Run("folder", func(t *testing.T) {
		out.Reset()
		folder, err := os.MkdirTemp("", "folder")
		require.NoError(t, err)
		defer os.RemoveAll(folder)
		app.RenameEventHandler(folder)

		var data map[string]interface{}
		err = json.Unmarshal(out.Bytes(), &data)
		require.NoError(t, err)

		level, ok := data["level"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "INFO", level)

		msg, ok := data["msg"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "receive: rename folder event", msg)

		event, ok := data["event"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "RENAME", event)

		path, ok := data["path"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, folder, path)
	})
}

func TestDeleteEventHandler(t *testing.T) {
	out := bytes.NewBuffer(nil)
	logger := testhelpers.NewTestLogger(t, out)
	app := New(1, 0, "", "", nil, logger)
	t.Run("file", func(t *testing.T) {
		file, err := os.CreateTemp("", "file")
		require.NoError(t, err)
		fpath := file.Name()
		file.Close()
		defer os.Remove(file.Name())
		app.DeleteEventHandler(fpath)

		var data map[string]interface{}
		err = json.Unmarshal(out.Bytes(), &data)
		require.NoError(t, err)

		level, ok := data["level"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "INFO", level)

		msg, ok := data["msg"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "receive: delete file event", msg)

		event, ok := data["event"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "DELETE", event)

		path, ok := data["path"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, fpath, path)
	})

	t.Run("folder", func(t *testing.T) {
		out.Reset()
		folder, err := os.MkdirTemp("", "folder")
		require.NoError(t, err)
		defer os.RemoveAll(folder)
		app.DeleteEventHandler(folder)

		var data map[string]interface{}
		err = json.Unmarshal(out.Bytes(), &data)
		require.NoError(t, err)

		level, ok := data["level"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "INFO", level)

		msg, ok := data["msg"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "receive: delete folder event", msg)

		event, ok := data["event"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "DELETE", event)

		path, ok := data["path"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, folder, path)
	})
}

func TestAttributeEventHandler(t *testing.T) {
	out := bytes.NewBuffer(nil)
	logger := testhelpers.NewTestLogger(t, out)
	app := New(1, 0, "", "", nil, logger)
	t.Run("file", func(t *testing.T) {
		file, err := os.CreateTemp("", "file")
		require.NoError(t, err)
		fpath := file.Name()
		file.Close()
		defer os.Remove(file.Name())
		app.AttributeEventHandler(fpath)

		var data map[string]interface{}
		err = json.Unmarshal(out.Bytes(), &data)
		require.NoError(t, err)

		level, ok := data["level"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "INFO", level)

		msg, ok := data["msg"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "receive: file attribute event", msg)

		event, ok := data["event"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "PERM", event)

		path, ok := data["path"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, fpath, path)
	})

	t.Run("folder", func(t *testing.T) {
		out.Reset()
		folder, err := os.MkdirTemp("", "folder")
		require.NoError(t, err)
		defer os.RemoveAll(folder)
		app.AttributeEventHandler(folder)

		var data map[string]interface{}
		err = json.Unmarshal(out.Bytes(), &data)
		require.NoError(t, err)

		level, ok := data["level"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "INFO", level)

		msg, ok := data["msg"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "receive: folder attribute event", msg)

		event, ok := data["event"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, "PERM", event)

		path, ok := data["path"].(string)
		require.Equal(t, true, ok)
		assert.Equal(t, folder, path)
	})
}
