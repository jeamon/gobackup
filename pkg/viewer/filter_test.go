package viewer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsEntryMatches(t *testing.T) {
	cases := []struct {
		name     string
		log      string
		date     string
		reg      string
		expected bool
	}{
		{"match exact path", `{"time":"2023-08-14T15:47:12.7081903Z", "path":".bak"}`, "2023-08-14", ".bak", true},
		{"does not match exact path", `{"time":"2023-08-15T15:47:12.7081903Z", "path":"file.bak"}`, "2023-08-15", ".bak", false},
		{"match all paths", `{"time":"2023-08-16T15:47:12.7081903Z", "path":"file.bak"}`, "2023-08-16", "*.*", true},
		{"does not match path", `{"time":"2023-08-17T15:47:12.7081903Z", "path":"file.txt"}`, "2023-08-17", "*.bak", false},
		{"match time and path", `{"time":"2023-08-18T15:47:12.7081903Z", "path":"_file_"}`, "2023-08-18", "*file*", true},
		{"invalid json", `{time:"2023-08-18T15:47:12.7081903Z", "path":"_file_"}`, "2023-08-18", "*file*", false},
		{"missing time json", `{"path":"_file_"}`, "2023-08-18", "*file*", false},
		{"missing path field", `{"time":"2023-08-18T15:47:12.7081903Z"}`, "2023-08-18", "*file*", false},
	}

	for _, tc := range cases {
		t.Run(tc.date, func(t *testing.T) {
			got := IsEntryMatches(tc.log, tc.date, tc.reg)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestFilter(t *testing.T) {
	file, err := os.CreateTemp("", "file")
	require.NoError(t, err)
	defer os.Remove(file.Name())
	_, err = file.WriteString(`{"time":"2023-08-14T15:47:12.7081903Z", "path":"file.zip"}`)
	require.NoError(t, err)
	_, err = file.WriteString(`{"time":"2023-08-15T15:47:12.7081903Z", "path":"file.bak"}`)
	require.NoError(t, err)
	defer file.Close()

	code, err := Filter(file, "2023-08-14", "file*")
	assert.NoError(t, err)
	assert.Equal(t, 0, code)
}
