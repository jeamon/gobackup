package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/jeamon/gobackup/pkg/events"
	"github.com/jeamon/gobackup/pkg/notifier"
	"github.com/jeamon/gobackup/pkg/testhelpers"
	"github.com/jeamon/gorsn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	w, err := notifier.New("./", nil)
	require.NoError(t, err)
	app := New(2, 0, "./", "dst", w, testhelpers.NewTestLogger(t, io.Discard))
	assert.Equal(t, 0, app.pid)
	assert.Equal(t, "./", app.srcFolder)
	assert.Equal(t, "dst", app.dstFolder)
	assert.Equal(t, w, app.notifier)
	assert.Equal(t, 2, cap(app.jobs))
}

func TestSigHandler(t *testing.T) {
	sigChan := make(chan os.Signal, 1)
	app := New(1, 0, "", "", nil, nil)
	go func() {
		time.Sleep(1 * time.Millisecond)
		sigChan <- syscall.SIGINT
	}()
	app.sigHandler(sigChan)
	open := true
	select {
	case _, open = <-app.stop:
	default:
	}
	assert.Equal(t, false, open)
}

func TestStart_Success(t *testing.T) {
	startCalled := false
	watcher := &testhelpers.MockMonitor{
		StartFunc: func(_ context.Context, quit <-chan struct{}, _ events.Queue) error {
			startCalled = true
			<-quit
			return nil
		},
	}

	src, err := os.MkdirTemp("", "source")
	require.NoError(t, err)
	defer os.RemoveAll(src)
	file, err := os.CreateTemp(src, "file")
	require.NoError(t, err)
	file.Close()

	// parent directory of backup folder.
	pfolder, err := os.MkdirTemp("", "folder")
	require.NoError(t, err)
	defer os.RemoveAll(pfolder)

	dst, err := os.MkdirTemp(pfolder, "backup")
	require.NoError(t, err)

	out := bytes.NewBuffer(nil)
	logger := testhelpers.NewTestLogger(t, out)
	app := New(1, 0, src, dst, watcher, logger)
	go func() {
		time.Sleep(1 * time.Second)
		app.Stop()
	}()
	code, err := app.start(1)
	assert.Equal(t, true, startCalled)
	assert.NoError(t, err)
	assert.Equal(t, 0, code)

	var data map[string]interface{}
	err = json.Unmarshal(out.Bytes(), &data)
	require.NoError(t, err)

	level, ok := data["level"].(string)
	require.Equal(t, true, ok)
	assert.Equal(t, "INFO", level)

	msg, ok := data["msg"].(string)
	require.Equal(t, true, ok)
	assert.Equal(t, "success: save backup folder state [success/fails: 0/0]", msg)

	event, ok := data["event"].(string)
	require.Equal(t, true, ok)
	assert.Equal(t, "SAVE", event)

	path, ok := data["path"].(string)
	require.Equal(t, true, ok)
	assert.Equal(t, pfolder, filepath.Dir(path))
	assert.Equal(t, ".zip", filepath.Ext(path))
}

func TestStart_Fail(t *testing.T) {
	watcher := &testhelpers.MockMonitor{
		StartFunc: func(_ context.Context, _ <-chan struct{}, _ events.Queue) error {
			return gorsn.ErrScanIsNotReady
		},
	}

	src, err := os.MkdirTemp("", "source")
	require.NoError(t, err)
	defer os.RemoveAll(src)
	file, err := os.CreateTemp(src, "file")
	require.NoError(t, err)
	file.Close()

	dst, err := os.MkdirTemp("", "backup")
	require.NoError(t, err)
	defer os.RemoveAll(dst)

	app := New(1, 0, src, dst, watcher, nil)
	code, err := app.start(1)
	assert.EqualError(t, err, fmt.Sprintf("failed to start files monitor: %v", gorsn.ErrScanIsNotReady))
	assert.Equal(t, 1, code)
}
