package repositories

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
)

// LanguageDetector detects whether a given directory contains a specific language.
type LanguageDetector interface {
	// DetectionFiles returns the ordered list of filenames/globs that
	// positively identify this language. First match wins.
	DetectionFiles() []string

	// Detect returns true if the given directory contains this language's markers.
	Detect(repoPath string) (bool, error)

	// Language returns the canonical Language entity for this detector.
	Language() entities.Language
}

// DetectWith checks if a language is present using the given FileChecker
// against the detector's DetectionFiles. This enables the same detection
// logic to work with both local filesystem and remote API-based file access.
func DetectWith(detector LanguageDetector, checker entities.FileChecker) (bool, error) {
	for _, f := range detector.DetectionFiles() {
		found, err := checker(f)
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}
	return false, nil
}
