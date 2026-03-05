package repositories

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
)

// LanguageProvider is the composite interface for a complete language implementation.
// It combines detection, version management, and dependency management.
type LanguageProvider interface {
	LanguageDetector
	VersionReader
	VersionWriter
	DependencyReader
	DependencyUpdater
}

// LanguageProviderWithValidation is a LanguageProvider that also supports build validation.
type LanguageProviderWithValidation interface {
	LanguageProvider
	BuildValidator
}

// LanguageInfo returns basic info about a language provider.
type LanguageInfo interface {
	// Language returns the canonical Language entity.
	Language() entities.Language
}
