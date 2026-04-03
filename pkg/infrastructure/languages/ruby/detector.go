package ruby

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects Ruby projects.
type Detector struct{}

// DetectionFiles returns the files that identify a Ruby project.
func (d *Detector) DetectionFiles() []string {
	return []string{"Gemfile", "*.gemspec"}
}

// Detect returns true if Gemfile or *.gemspec exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	return repositories.DetectWith(d, fileutil.LocalFileChecker(repoPath))
}

// Language returns the Ruby language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguageRuby
}
