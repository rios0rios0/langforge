package golang

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects Go projects by the presence of go.mod.
type Detector struct{}

// DetectionFiles returns the files that identify a Go project.
func (d *Detector) DetectionFiles() []string {
	return []string{"go.mod"}
}

// Detect returns true if go.mod exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	return repositories.DetectWith(d, fileutil.LocalFileChecker(repoPath))
}

// Language returns the Go language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguageGo
}
