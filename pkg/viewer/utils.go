package viewer

import (
	"fmt"
	"io"
	"os"
	"time"
)

// IsValidFilters checks if the `date` and `regex` values
// provided as program arguments are valid ones.
func IsValidFilters(date, regex string) bool {
	if len(regex) == 0 {
		return false
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return false
	}
	return true
}

// Print writes a log message to the provided output.
func Print(out io.Writer, data string) {
	fmt.Fprintln(out, data)
}

// Open opens in read only mode a file.
func Open(file string) (*os.File, error) {
	f, err := os.Open(file)
	return f, err
}
