package java_maven

import (
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects Java/Maven projects by the presence of pom.xml.
type Detector struct{}

// DetectionFiles returns the files that identify a Maven project.
func (d *Detector) DetectionFiles() []string {
	return []string{"pom.xml"}
}

// Detect returns true if pom.xml exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	return fileutil.Exists(filepath.Join(repoPath, "pom.xml")), nil
}

// Language returns the Java/Maven language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguageJavaMaven
}
