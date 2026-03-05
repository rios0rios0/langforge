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
		if strings.ContainsAny(pathOrPattern, "*?[") {
			matches, err := GlobFiles(repoPath, pathOrPattern)
			if err != nil {
				return false, err
			}
			return len(matches) > 0, nil
		}
		return Exists(filepath.Join(repoPath, pathOrPattern)), nil
	}
}
