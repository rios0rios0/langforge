package dart

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects Dart and Flutter projects by the presence of pubspec.yaml.
type Detector struct{}

// DetectionFiles returns the files that identify a Dart project.
// Flutter shares the same manifest, so a single marker covers both.
func (d *Detector) DetectionFiles() []string {
	return []string{Pubspec}
}

// Detect returns true if pubspec.yaml exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	return repositories.DetectWith(d, fileutil.LocalFileChecker(repoPath))
}

// Language returns the Dart language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguageDart
}
