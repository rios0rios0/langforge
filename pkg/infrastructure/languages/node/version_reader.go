package node

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// VersionReader reads the version from package.json.
type VersionReader struct{}

// VersionFiles returns the files inspected for the version.
func (r *VersionReader) VersionFiles() []string {
	return []string{"package.json"}
}

// ReadVersion reads the version field from package.json.
func (r *VersionReader) ReadVersion(repoPath string) (entities.Version, error) {
	content, err := fileutil.ReadFile(filepath.Join(repoPath, "package.json"))
	if err != nil {
		return entities.Version{}, fmt.Errorf("reading package.json: %w", err)
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal([]byte(content), &pkg); err != nil {
		return entities.Version{}, fmt.Errorf("parsing package.json: %w", err)
	}
	if pkg.Version == "" {
		return entities.Version{}, fmt.Errorf("no version field in package.json")
	}
	return entities.NewVersion(pkg.Version)
}
