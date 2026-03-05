package entitybuilders

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
)

// LanguageProviderStub is a test double for LanguageProvider.
type LanguageProviderStub struct {
	lang           entities.Language
	detectionFiles []string
	detectResult   bool
	detectErr      error
	version        entities.Version
	versionErr     error
	versionFiles   []string
	filesChanged   []string
	writeErr       error
	deps           []entities.Dependency
	depsErr        error
	updateErr      error
	commands       []string
}

// DetectionFiles returns the configured detection files.
func (s *LanguageProviderStub) DetectionFiles() []string { return s.detectionFiles }

// Detect returns the configured detection result.
func (s *LanguageProviderStub) Detect(_ string) (bool, error) { return s.detectResult, s.detectErr }

// Language returns the configured language.
func (s *LanguageProviderStub) Language() entities.Language { return s.lang }

// VersionFiles returns the configured version files.
func (s *LanguageProviderStub) VersionFiles() []string { return s.versionFiles }

// ReadVersion returns the configured version.
func (s *LanguageProviderStub) ReadVersion(_ string) (entities.Version, error) {
	return s.version, s.versionErr
}

// WriteVersion returns the configured write error.
func (s *LanguageProviderStub) WriteVersion(_ string, _ entities.Version) error {
	return s.writeErr
}

// FilesChanged returns the configured files changed.
func (s *LanguageProviderStub) FilesChanged(_ string) ([]string, error) {
	return s.filesChanged, nil
}

// ReadDependencies returns the configured dependencies.
func (s *LanguageProviderStub) ReadDependencies(_ string) ([]entities.Dependency, error) {
	return s.deps, s.depsErr
}

// UpdateAll returns the configured update error.
func (s *LanguageProviderStub) UpdateAll(_ string) error { return s.updateErr }

// Commands returns the configured commands.
func (s *LanguageProviderStub) Commands() []string { return s.commands }

// LanguageProviderStubBuilder builds LanguageProviderStub instances.
type LanguageProviderStubBuilder struct {
	stub *LanguageProviderStub
}

// NewLanguageProviderStubBuilder creates a new builder with defaults.
func NewLanguageProviderStubBuilder() *LanguageProviderStubBuilder {
	return &LanguageProviderStubBuilder{
		stub: &LanguageProviderStub{
			lang:           entities.LanguageUnknown,
			detectionFiles: []string{},
			versionFiles:   []string{},
			filesChanged:   []string{},
			commands:       []string{},
		},
	}
}

// WithLanguage sets the language.
func (b *LanguageProviderStubBuilder) WithLanguage(lang entities.Language) *LanguageProviderStubBuilder {
	b.stub.lang = lang
	return b
}

// WithDetectionFiles sets the detection files.
func (b *LanguageProviderStubBuilder) WithDetectionFiles(files []string) *LanguageProviderStubBuilder {
	b.stub.detectionFiles = files
	return b
}

// WithDetectResult sets the detect result.
func (b *LanguageProviderStubBuilder) WithDetectResult(result bool, err error) *LanguageProviderStubBuilder {
	b.stub.detectResult = result
	b.stub.detectErr = err
	return b
}

// WithVersion sets the version to return.
func (b *LanguageProviderStubBuilder) WithVersion(v entities.Version, err error) *LanguageProviderStubBuilder {
	b.stub.version = v
	b.stub.versionErr = err
	return b
}

// WithWriteError sets the write error.
func (b *LanguageProviderStubBuilder) WithWriteError(err error) *LanguageProviderStubBuilder {
	b.stub.writeErr = err
	return b
}

// WithDependencies sets the dependencies.
func (b *LanguageProviderStubBuilder) WithDependencies(deps []entities.Dependency, err error) *LanguageProviderStubBuilder {
	b.stub.deps = deps
	b.stub.depsErr = err
	return b
}

// WithUpdateError sets the update error.
func (b *LanguageProviderStubBuilder) WithUpdateError(err error) *LanguageProviderStubBuilder {
	b.stub.updateErr = err
	return b
}

// WithCommands sets the commands.
func (b *LanguageProviderStubBuilder) WithCommands(commands []string) *LanguageProviderStubBuilder {
	b.stub.commands = commands
	return b
}

// Build builds the stub.
func (b *LanguageProviderStubBuilder) Build() *LanguageProviderStub {
	return b.stub
}
