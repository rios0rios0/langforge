package node

import (
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects Node/TypeScript projects by the presence of package.json.
type Detector struct{}

// DetectionFiles returns the files that identify a Node project.
func (d *Detector) DetectionFiles() []string {
	return []string{"package.json", "tsconfig.json"}
}

// Detect returns true if package.json exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	return fileutil.Exists(filepath.Join(repoPath, "package.json")), nil
}

// Language returns the Node language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguageNode
}
