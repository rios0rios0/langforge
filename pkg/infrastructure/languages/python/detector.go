package python

import (
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects Python projects.
type Detector struct{}

// DetectionFiles returns the files that identify a Python project.
func (d *Detector) DetectionFiles() []string {
	return []string{"pyproject.toml", "setup.py"}
}

// Detect returns true if pyproject.toml or setup.py exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	for _, f := range d.DetectionFiles() {
		if fileutil.Exists(filepath.Join(repoPath, f)) {
			return true, nil
		}
	}
	return false, nil
}

// Language returns the Python language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguagePython
}
