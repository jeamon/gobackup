package logger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupLogger(t *testing.T) {
	folder, err := os.MkdirTemp("", "folder")
	require.NoError(t, err)
	defer os.RemoveAll(folder)

	file, err := os.CreateTemp(folder, "file")
	require.NoError(t, err)
	filePath := file.Name()
	file.Close()

	t.Run("should open log file", func(t *testing.T) {
		logfile, logger, err := New(filePath, "commit", "tag", 0)
		require.NoError(t, err)
		require.FileExists(t, filePath)
		assert.NotNil(t, logfile)
		assert.NotNil(t, logger)
		logfile.Close()
	})
	t.Run("should create log file", func(t *testing.T) {
		filePath := filepath.Join(folder, "log.file")
		logfile, logger, err := New(filePath, "commit", "tag", 0)
		require.NoError(t, err)
		require.FileExists(t, filePath)
		assert.NotNil(t, logfile)
		assert.NotNil(t, logger)
		logfile.Close()
	})

	t.Run("should fail", func(t *testing.T) {
		filePath := filepath.Join(folder, "noexist.folder", "log.file")
		logfile, logger, err := New(filePath, "commit", "tag", 0)
		require.Error(t, err)
		require.NoFileExists(t, filePath)
		assert.Nil(t, logfile)
		assert.Nil(t, logger)
	})
}
