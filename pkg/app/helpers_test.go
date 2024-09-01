package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateBackupFileContent(t *testing.T) {
	sfile, err := os.CreateTemp("", "gobackup")
	require.NoError(t, err)
	sfilepath := sfile.Name()
	defer os.Remove(sfilepath)
	_, err = sfile.WriteString("original content")
	require.NoError(t, err)
	sfile.Close()

	app := &App{dstFolder: filepath.Dir(sfilepath)}
	err = app.UpdateBackupFileContent(sfilepath)
	require.NoError(t, err)
	backupFilepath := sfilepath + backupFileExtension
	defer os.Remove(backupFilepath)

	data, err := os.ReadFile(backupFilepath)
	require.NoError(t, err)
	ok := bytes.Equal([]byte("original content"), data)
	assert.Equal(t, true, ok)
}

func TestCreateBackupFile(t *testing.T) {
	sfile, err := os.CreateTemp("", "gobackup")
	require.NoError(t, err)
	sfilepath := sfile.Name()
	defer os.Remove(sfilepath)
	sfile.Close()

	app := &App{dstFolder: filepath.Dir(sfilepath)}
	err = app.CreateBackupFile(sfilepath)
	require.NoError(t, err)

	backupFilepath := sfilepath + backupFileExtension
	fi, err := os.Stat(backupFilepath)
	require.NoError(t, err)
	assert.Equal(t, true, fi.Mode().IsRegular())
	os.Remove(backupFilepath)
}

func TestIsImmediateDelete(t *testing.T) {
	cases := []struct {
		path     string
		expected bool
	}{
		{"", false},
		{"./delete_", false},
		{"delete__", true},
		{"_delete_file", false},
		{" delete_file", false},
		{"delete_file", true},
		{"Delete_file", false},
		{"delete__file", true},
		{"delete_", false},
		{"delete_ ", true},
	}
	app := &App{}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			got := app.IsImmediateDelete(tc.path)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestIsScheduleDelete(t *testing.T) {
	src, err := os.MkdirTemp("", "source")
	require.NoError(t, err)
	defer os.RemoveAll(src)
	dst, err := os.MkdirTemp("", "backup")
	require.NoError(t, err)
	defer os.RemoveAll(dst)
	emptyTime := time.Time{}
	path := func(folder, filename string) string {
		return filepath.Join(folder, filename)
	}
	maketime := func(ts string) time.Time {
		t, err := time.Parse(time.RFC3339, ts)
		if err == nil {
			return t
		}
		return emptyTime
	}

	cases := []struct {
		path     string
		src      string
		dst      string
		datetime time.Time
		match    bool
	}{
		{path(src, ""), "", "", emptyTime, false},
		{path(src, "delete_filename"), "", "", emptyTime, false},
		{path(src, "delete_2023-08-14"), "", "", emptyTime, false},
		{path(src, "delete_2023-08-14_filename"), "", "", emptyTime, false},
		{path(src, "delete_2023-08-14T00:00:00Z"), "", "", emptyTime, false},
		{path(src, "delete_2023-08-14T00:00:00Z_"), "", "", emptyTime, false},
		{
			path(src, "delete_2023-08-14T00:00:00-02:00_ _"),
			path(src, " _"),
			path(dst, " _"+backupFileExtension),
			maketime("2023-08-14T00:00:00-02:00"),
			true,
		},
		{
			path(src, "delete_2023-08-14T00:00:00Z__filename"),
			path(src, "_filename"),
			path(dst, "_filename"+".bak"),
			maketime("2023-08-14T00:00:00Z"),
			true,
		},
		{
			path(src, "delete_2023-08-14T00:00:00-02:00__filename"),
			path(src, "_filename"),
			path(dst, "_filename"+".bak"),
			maketime("2023-08-14T00:00:00-02:00"),
			true,
		},
	}

	app := &App{srcFolder: src, dstFolder: dst}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			s, d, dt, m := app.IsScheduleDelete(tc.path)
			assert.Equal(t, tc.match, m)
			assert.Equal(t, tc.src, s)
			assert.Equal(t, tc.dst, d)
			assert.Equal(t, tc.datetime, dt)
		})
	}
}
