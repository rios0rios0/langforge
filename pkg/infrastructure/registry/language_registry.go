package registry

import (
	"fmt"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
)

// LanguageRegistry maps Language → LanguageProvider.
// Populated via Register() at init() time in each language package.
type LanguageRegistry struct {
	providers []repositories.LanguageProvider
	byLang    map[entities.Language]repositories.LanguageProvider
}

// NewLanguageRegistry creates a new empty LanguageRegistry.
func NewLanguageRegistry() *LanguageRegistry {
	return &LanguageRegistry{
		byLang: make(map[entities.Language]repositories.LanguageProvider),
	}
}

// Register adds a LanguageProvider to the registry.
func (r *LanguageRegistry) Register(provider repositories.LanguageProvider) {
	r.providers = append(r.providers, provider)
	r.byLang[provider.Language()] = provider
}

// Detect scans the given repo path and returns the first matching LanguageProvider.
// Returns an error if no language is detected.
func (r *LanguageRegistry) Detect(repoPath string) (repositories.LanguageProvider, error) {
	for _, p := range r.providers {
		matched, err := p.Detect(repoPath)
		if err != nil {
			return nil, fmt.Errorf("detection error for %s: %w", p.Language(), err)
		}
		if matched {
			return p, nil
		}
	}
	return nil, fmt.Errorf("no supported language detected in %q", repoPath)
}

// Get returns the LanguageProvider for the given Language.
// Returns an error if no provider is registered for that language.
func (r *LanguageRegistry) Get(lang entities.Language) (repositories.LanguageProvider, error) {
	p, ok := r.byLang[lang]
	if !ok {
		return nil, fmt.Errorf("no provider registered for language %q", lang)
	}
	return p, nil
}

// DetectWithChecker scans registered providers using the given FileChecker
// and returns the first matching LanguageProvider. This enables detection
// against remote APIs or other non-filesystem sources.
func (r *LanguageRegistry) DetectWithChecker(checker entities.FileChecker) (repositories.LanguageProvider, error) {
	for _, p := range r.providers {
		matched, err := repositories.DetectWith(p, checker)
		if err != nil {
			return nil, fmt.Errorf("detection error for %s: %w", p.Language(), err)
		}
		if matched {
			return p, nil
		}
	}
	return nil, fmt.Errorf("no supported language detected")
}

// DetectAllWithChecker scans registered providers using the given FileChecker
// and returns all matching LanguageProviders.
func (r *LanguageRegistry) DetectAllWithChecker(checker entities.FileChecker) ([]repositories.LanguageProvider, error) {
	var matched []repositories.LanguageProvider
	for _, p := range r.providers {
		ok, err := repositories.DetectWith(p, checker)
		if err != nil {
			return nil, fmt.Errorf("detection error for %s: %w", p.Language(), err)
		}
		if ok {
			matched = append(matched, p)
		}
	}
	return matched, nil
}

// Languages returns the list of registered language names.
func (r *LanguageRegistry) Languages() []entities.Language {
	langs := make([]entities.Language, 0, len(r.byLang))
	for lang := range r.byLang {
		langs = append(langs, lang)
	}
	return langs
}
