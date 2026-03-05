package node

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// DependencyReader reads dependencies from package.json.
type DependencyReader struct{}

// ReadDependencies parses package.json and returns the list of dependencies.
func (r *DependencyReader) ReadDependencies(repoPath string) ([]entities.Dependency, error) {
	content, err := fileutil.ReadFile(filepath.Join(repoPath, "package.json"))
	if err != nil {
		return nil, fmt.Errorf("reading package.json: %w", err)
	}
	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err = json.Unmarshal([]byte(content), &pkg); err != nil {
		return nil, fmt.Errorf("parsing package.json: %w", err)
	}
	var deps []entities.Dependency
	for name, version := range pkg.Dependencies {
		deps = append(deps, entities.NewDependency(name, version, "", "package.json"))
	}
	for name, version := range pkg.DevDependencies {
		deps = append(deps, entities.NewDependency(name, version, "", "package.json"))
	}
	return deps, nil
}
