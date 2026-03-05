package repositories

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
)

// VersionWriter writes a new version to the appropriate files in a project directory.
type VersionWriter interface {
	// WriteVersion updates the version in the appropriate files.
	WriteVersion(repoPath string, version entities.Version) error

	// FilesChanged returns the list of files that will be modified.
	FilesChanged(repoPath string) ([]string, error)
}
