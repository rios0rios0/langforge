package repositories

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
)

// DependencyReader reads dependencies from a project directory.
type DependencyReader interface {
	// ReadDependencies reads all dependencies from the given repo path.
	ReadDependencies(repoPath string) ([]entities.Dependency, error)
}
