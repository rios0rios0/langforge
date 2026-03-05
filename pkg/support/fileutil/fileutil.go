package fileutil

import (
	"os"
	"path/filepath"
)

// Exists returns true if the path exists and is accessible.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ReadFile reads the content of a file and returns it as a string.
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- path is constructed from repo root
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes content to a file, creating it if necessary.
func WriteFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o600) // #nosec G306
}

// GlobFiles returns files matching a pattern rooted at dir.
func GlobFiles(dir, pattern string) ([]string, error) {
	return filepath.Glob(filepath.Join(dir, pattern))
}

// FindFirst returns the first file in dir that matches any of the given names.
// Returns "" if none found.
func FindFirst(dir string, names []string) string {
	for _, name := range names {
		if Exists(filepath.Join(dir, name)) {
			return filepath.Join(dir, name)
		}
	}
	return ""
}
