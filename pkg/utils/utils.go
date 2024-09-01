package utils

import "os"

// IsDirPath checks if provided path is an accessible directory.
func IsDirPath(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// CreateFolder creates complete absolute path.
func CreateFolder(path string) error {
	return os.MkdirAll(path, 0o644)
}

// DeleteFile removes a file based on its absolute path.
func DeleteFile(path string) error {
	return os.Remove(path)
}
