package python

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects Python projects.
type Detector struct{}

// DetectionFiles returns the files that identify a Python project.
func (d *Detector) DetectionFiles() []string {
	return []string{"pyproject.toml", "setup.py", "requirements.txt"}
}

// Detect returns true if pyproject.toml, setup.py, or requirements.txt exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	return repositories.DetectWith(d, fileutil.LocalFileChecker(repoPath))
}

// Language returns the Python language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguagePython
}
