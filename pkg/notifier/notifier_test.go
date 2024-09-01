package notifier

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jeamon/gobackup/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddFolder(t *testing.T) {
	src, err := os.MkdirTemp("", "source")
	require.NoError(t, err)
	defer os.RemoveAll(src)

	f1, err := os.CreateTemp(src, "file")
	require.NoError(t, err)
	f1.Close()

	folder, err := os.MkdirTemp(src, "folder")
	require.NoError(t, err)

	f2, err := os.CreateTemp(folder, "file")
	require.NoError(t, err)
	f2.Close()

	_, err = New(src, nil)
	require.NoError(t, err)
}

func TestStart(t *testing.T) {
	jobs := make(events.Queue, 1)
	quit := make(chan struct{})

	src, err := os.MkdirTemp("", "source")
	require.NoError(t, err)
	defer os.RemoveAll(src)
	notifier, err := New(src, nil)
	require.NoError(t, err)
	go func() {
		file, err := os.Create(filepath.Join(src, "file"))
		require.NoError(t, err)
		require.NoError(t, file.Close())
		time.Sleep(3 * time.Second)
		quit <- struct{}{}
	}()

	notifier.Start(context.Background(), quit, jobs)
	select {
	case ce := <-jobs:
		assert.Equal(t, filepath.Join(src, "file"), ce.Path)
		assert.Equal(t, events.CREATE, ce.Ops)
	case <-time.After(5 * time.Second):
		t.Error("failed because taking too much time")
	}
}
