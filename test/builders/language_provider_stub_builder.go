package builders

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/test/doubles"
	testkit "github.com/rios0rios0/testkit/pkg/test"
)

// LanguageProviderStubBuilder builds LanguageProviderStub instances using the builder pattern.
type LanguageProviderStubBuilder struct {
	*testkit.BaseBuilder

	lang            entities.Language
	detectionFiles  []string
	detectResult    bool
	detectErr       error
	version         entities.Version
	versionErr      error
	versionFiles    []string
	filesChanged    []string
	filesChangedErr error
	writeErr        error
	deps            []entities.Dependency
	depsErr         error
	updateErr       error
	commands        []string
}

// NewLanguageProviderStubBuilder creates a new builder with defaults.
func NewLanguageProviderStubBuilder() *LanguageProviderStubBuilder {
	return &LanguageProviderStubBuilder{
		BaseBuilder:    testkit.NewBaseBuilder(),
		lang:           entities.LanguageUnknown,
		detectionFiles: []string{},
		versionFiles:   []string{},
		filesChanged:   []string{},
		commands:       []string{},
	}
}

func (b *LanguageProviderStubBuilder) WithLanguage(lang entities.Language) *LanguageProviderStubBuilder {
	b.lang = lang
	return b
}

func (b *LanguageProviderStubBuilder) WithDetectionFiles(files []string) *LanguageProviderStubBuilder {
	b.detectionFiles = files
	return b
}

func (b *LanguageProviderStubBuilder) WithDetectResult(
	result bool, err error,
) *LanguageProviderStubBuilder {
	b.detectResult = result
	b.detectErr = err
	return b
}

func (b *LanguageProviderStubBuilder) WithVersion(
	v entities.Version, err error,
) *LanguageProviderStubBuilder {
	b.version = v
	b.versionErr = err
	return b
}

func (b *LanguageProviderStubBuilder) WithWriteError(err error) *LanguageProviderStubBuilder {
	b.writeErr = err
	return b
}

func (b *LanguageProviderStubBuilder) WithFilesChangedError(err error) *LanguageProviderStubBuilder {
	b.filesChangedErr = err
	return b
}

func (b *LanguageProviderStubBuilder) WithDependencies(
	deps []entities.Dependency, err error,
) *LanguageProviderStubBuilder {
	b.deps = deps
	b.depsErr = err
	return b
}

func (b *LanguageProviderStubBuilder) WithUpdateError(err error) *LanguageProviderStubBuilder {
	b.updateErr = err
	return b
}

func (b *LanguageProviderStubBuilder) WithCommands(commands []string) *LanguageProviderStubBuilder {
	b.commands = commands
	return b
}

func (b *LanguageProviderStubBuilder) Build() any {
	return &doubles.LanguageProviderStub{
		LangValue:           b.lang,
		DetectionFilesValue: b.detectionFiles,
		DetectResultValue:   b.detectResult,
		DetectErrValue:      b.detectErr,
		VersionValue:        b.version,
		VersionErrValue:     b.versionErr,
		VersionFilesValue:   b.versionFiles,
		FilesChangedValue:   b.filesChanged,
		FilesChangedErr:     b.filesChangedErr,
		WriteErrValue:       b.writeErr,
		DepsValue:           b.deps,
		DepsErrValue:        b.depsErr,
		UpdateErrValue:      b.updateErr,
		CommandsValue:       b.commands,
	}
}
