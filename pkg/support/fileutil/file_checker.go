package fileutil

import (
	"path/filepath"
	"strings"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
)

// LocalFileChecker returns a FileChecker that uses the local filesystem.
// For glob patterns (containing *, ?, or [), it uses [filepath.Glob].
// For exact paths, it checks file existence with [os.Stat].
func LocalFileChecker(repoPath string) entities.FileChecker {
	return func(pathOrPattern string) (bool, error) {
		if IsGlobPattern(pathOrPattern) {
			matches, err := GlobFiles(repoPath, pathOrPattern)
			if err != nil {
				return false, err
			}
			return len(matches) > 0, nil
		}
		return Exists(filepath.Join(repoPath, pathOrPattern)), nil
	}
}

// NewFileChecker builds a FileChecker from two callback functions:
// exactCheck is called for exact file paths, and globCheck is called
// for glob patterns. This enables consumers to plug in remote
// API-based file access (e.g. via gitforge's FileAccessProvider)
// without langforge needing to import provider-specific types.
//
// Exported for use by autoupdate and autobump to build remote-compatible
// FileCheckers that query Git hosting APIs instead of the local filesystem.
func NewFileChecker(
	exactCheck func(path string) (bool, error),
	globCheck func(pattern string) (bool, error),
) entities.FileChecker {
	return func(pathOrPattern string) (bool, error) {
		if IsGlobPattern(pathOrPattern) {
			return globCheck(pathOrPattern)
		}
		return exactCheck(pathOrPattern)
	}
}

// IsGlobPattern returns true if the path contains glob metacharacters (*, ?, [).
func IsGlobPattern(path string) bool {
	return strings.ContainsAny(path, "*?[")
}

// ExtractExtension extracts the file extension from a glob pattern.
// For example, "*.tf" returns ".tf" and "*.hcl" returns ".hcl".
// Returns the input unchanged if no dot is found.
func ExtractExtension(pattern string) string {
	if idx := strings.LastIndex(pattern, "."); idx >= 0 {
		return pattern[idx:]
	}
	return pattern
}
