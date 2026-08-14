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

// CompositeProvider assembles a LanguageProviderFull out of the parts a
// language package implements. Every ecosystem composes the same seven ports in
// the same way, so they share this one type instead of each declaring its own
// identical struct.
//
// A language that implements only some of the ports leaves the rest nil, and
// the composite still satisfies the narrower interface (LanguageProvider or
// LanguageProviderWithValidation) its callers ask for. Only the ports it was
// given may be called: the assertion below is that the type can stand in for
// LanguageProviderFull, not that every instance carries a full set.
type CompositeProvider struct {
	LanguageDetector
	VersionReader
	VersionWriter
	DependencyReader
	DependencyUpdater
	BuildValidator
	RuntimeManager
}

var _ LanguageProviderFull = (*CompositeProvider)(nil)

// FilesChanged returns the union of the files the VersionWriter and the
// DependencyUpdater each report, writer's first and without duplicates.
//
// Both embedded ports declare FilesChanged, which makes the promoted method
// ambiguous, so the composite has to resolve it explicitly. Merging is the
// resolution every language wants: the caller asks which files an update
// touches, and both halves of the answer matter.
func (p *CompositeProvider) FilesChanged(repoPath string) ([]string, error) {
	files, err := p.VersionWriter.FilesChanged(repoPath)
	if err != nil {
		return nil, err
	}
	updated, err := p.DependencyUpdater.FilesChanged(repoPath)
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
