package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewLogs(t *testing.T) {
	t.Run("invalid filters", func(t *testing.T) {
		code, err := ViewLogs("logfile", "23-08-22", "")
		assert.Equal(t, 1, code)
		assert.EqualError(t, err, "invalid date and/or regex")
	})

	t.Run("log file does not exist", func(t *testing.T) {
		code, err := ViewLogs("logfile", "2023-08-22", "*.zip")
		assert.Equal(t, 1, code)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
	t.Run("run logs filter", func(t *testing.T) {
		file, err := os.CreateTemp("", "log")
		require.NoError(t, err)
		defer os.Remove(file.Name())
		file.Close()
		code, err := ViewLogs(file.Name(), "2023-08-22", "*.zip")
		assert.Equal(t, 0, code)
		assert.NoError(t, err)
	})
}
