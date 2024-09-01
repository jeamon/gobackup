package viewer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// IsEntryMatches returns wether the timestamp and filename specified
// into a given log entry matches the date and regex respectively.
// It decodes the log entry into a map so that each log entry line
// could have any set of fields.
func IsEntryMatches(logEntry, date, reg string) bool {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(logEntry), &data)
	if err != nil {
		return false
	}
	ts, ok := data["time"]
	if !ok || !strings.HasPrefix(ts.(string), date) {
		return false
	}
	path, ok := data["path"]
	if !ok {
		return false
	}
	if match, err := filepath.Match(reg, filepath.Base(path.(string))); err != nil || !match {
		return false
	}

	return true
}

// Filter process the content defined into `file` variable and displays
// all entries produced at the date `date` involving the filename matching
// the regex `reg`. The log file must be into the same folder as the program.
func Filter(file io.Reader, date, reg string) (int, error) {
	scanner := bufio.NewScanner(file)
	var logEntry string
	var match bool

	for scanner.Scan() {
		logEntry = strings.TrimSpace(scanner.Text())
		if logEntry == "" {
			continue
		}
		match = IsEntryMatches(logEntry, date, reg)
		if match {
			Print(os.Stdout, logEntry)
		}
	}

	if err := scanner.Err(); err != nil {
		return 1, fmt.Errorf("failed to filter file: %v", err)
	}
	return 0, nil
}
