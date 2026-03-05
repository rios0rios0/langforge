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
