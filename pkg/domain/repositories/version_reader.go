package repositories

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
)

// VersionReader reads the canonical version from a project directory.
type VersionReader interface {
	// ReadVersion reads the canonical version from the given repo path.
	ReadVersion(repoPath string) (entities.Version, error)

	// VersionFiles returns which files this reader inspects.
	VersionFiles() []string
}
