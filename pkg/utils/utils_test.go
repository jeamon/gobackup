package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsDirPath(t *testing.T) {
	dname, err := os.MkdirTemp("", "backup")
	require.NoError(t, err)
	defer os.RemoveAll(dname)

	file, err := os.CreateTemp("", "file")
	require.NoError(t, err)
	file.Close()
	defer os.Remove(file.Name())

	cases := []struct {
		path     string
		expected bool
	}{
		{"", false},
		{"./", true},
		{dname, true},
		{file.Name(), false},
		{dname + "noexist", false},
	}
	for _, tc := range cases {
		got := IsDirPath(tc.path)
		assert.Equal(t, tc.expected, got)
	}
}

func TestDeleteFile(t *testing.T) {
	sfile, err := os.CreateTemp("", "gobackup")
	require.NoError(t, err)
	sfilepath := sfile.Name()
	defer os.Remove(sfilepath)
	sfile.Close()

	err = DeleteFile(sfilepath)
	require.NoError(t, err)
	_, err = os.Stat(sfilepath)
	assert.Equal(t, true, os.IsNotExist(err))
}
