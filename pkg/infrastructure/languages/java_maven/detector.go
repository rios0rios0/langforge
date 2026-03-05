package java_maven

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
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
	return repositories.DetectWith(d, fileutil.LocalFileChecker(repoPath))
}

// Language returns the Java/Maven language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguageJavaMaven
}
