package javagradle

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects Java/Gradle projects.
type Detector struct{}

// DetectionFiles returns the files that identify a Gradle project.
func (d *Detector) DetectionFiles() []string {
	return []string{"build.gradle", "build.gradle.kts", "settings.gradle", "settings.gradle.kts"}
}

// Detect returns true if any Gradle marker file exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	return repositories.DetectWith(d, fileutil.LocalFileChecker(repoPath))
}

// Language returns the Java/Gradle language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguageJavaGradle
}
