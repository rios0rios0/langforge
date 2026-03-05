package java_gradle

import (
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects Java/Gradle projects.
type Detector struct{}

// DetectionFiles returns the files that identify a Gradle project.
func (d *Detector) DetectionFiles() []string {
	return []string{"build.gradle", "build.gradle.kts", "settings.gradle", "settings.gradle.kts"}
}

// Detect returns true if build.gradle or build.gradle.kts exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	for _, f := range []string{"build.gradle", "build.gradle.kts"} {
		if fileutil.Exists(filepath.Join(repoPath, f)) {
			return true, nil
		}
	}
	return false, nil
}

// Language returns the Java/Gradle language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguageJavaGradle
}
