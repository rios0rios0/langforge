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

// LanguageProviderFull is a LanguageProvider with build validation and runtime management.
type LanguageProviderFull interface {
	LanguageProviderWithValidation
	RuntimeManager
}

// LanguageInfo returns basic info about a language provider.
type LanguageInfo interface {
	// Language returns the canonical Language entity.
	Language() entities.Language
}

// MergeFilesChanged returns the union of the files a VersionWriter and a
// DependencyUpdater each report, writer's first and without duplicates.
//
// Every composite provider embeds both, and both declare FilesChanged, so the
// method is ambiguous on the embedded structs and each Provider has to define
// it explicitly. This is that definition: providers delegate here instead of
// each carrying its own copy of the merge.
func MergeFilesChanged(writer VersionWriter, updater DependencyUpdater, repoPath string) ([]string, error) {
	files, err := writer.FilesChanged(repoPath)
	if err != nil {
		return nil, err
	}
	updated, err := updater.FilesChanged(repoPath)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{}, len(files))
	for _, f := range files {
		seen[f] = struct{}{}
	}
	for _, f := range updated {
		if _, ok := seen[f]; !ok {
			files = append(files, f)
		}
	}
	return files, nil
}
