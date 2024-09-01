package viewer

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsValidFilters(t *testing.T) {
	cases := []struct {
		date     string
		reg      string
		expected bool
	}{
		{"", "", false},
		{"2023-08-14", "", false},
		{"2023-14-08", "*.bak", false},
		{"2023:08:08", "*.*", false},
		{"2023-08-14 14:00:00", "*.*", false},
		{"2023-08-14", "*.*", true},
		{"2023-08-14", "delete_*", true},
	}

	for _, tc := range cases {
		t.Run(tc.date+"/"+tc.reg, func(t *testing.T) {
			got := IsValidFilters(tc.date, tc.reg)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestPrint(t *testing.T) {
	out := bytes.NewBuffer(nil)
	data := "content to print"
	Print(out, data)
	assert.Equal(t, data+"\n", out.String())
}

func TestOpen(t *testing.T) {
	file, err := os.CreateTemp("", "log.file")
	require.NoError(t, err)
	file.Close()

	f, err := Open(file.Name())
	require.NoError(t, err)
	f.Close()
}
