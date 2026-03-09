package dockerfile

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects container projects by the presence of Dockerfile files.
type Detector struct{}

// DetectionFiles returns the files that identify a container project.
func (d *Detector) DetectionFiles() []string {
	return []string{
		"Dockerfile",
		"Dockerfile.*",
		"*.Dockerfile",
	}
}

// Detect returns true if any Dockerfile exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	return repositories.DetectWith(d, fileutil.LocalFileChecker(repoPath))
}

// Language returns the Dockerfile language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguageDockerfile
}
