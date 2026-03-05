package node

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// VersionWriter updates the version in package.json.
type VersionWriter struct{}

// FilesChanged returns the list of files that will be modified.
func (w *VersionWriter) FilesChanged(repoPath string) ([]string, error) {
	return []string{filepath.Join(repoPath, "package.json")}, nil
}

// WriteVersion updates the version field in package.json.
func (w *VersionWriter) WriteVersion(repoPath string, version entities.Version) error {
	path := filepath.Join(repoPath, "package.json")
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading package.json: %w", err)
	}
	var raw map[string]json.RawMessage
	if err = json.Unmarshal([]byte(content), &raw); err != nil {
		return fmt.Errorf("parsing package.json: %w", err)
	}
	versionJSON, err := json.Marshal(version.String())
	if err != nil {
		return fmt.Errorf("marshaling version: %w", err)
	}
	raw["version"] = versionJSON
	out, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling package.json: %w", err)
	}
	return fileutil.WriteFile(path, string(out)+"\n")
}
