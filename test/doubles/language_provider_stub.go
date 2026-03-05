package doubles

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
)

// LanguageProviderStub is a test double for LanguageProvider.
type LanguageProviderStub struct {
	LangValue           entities.Language
	DetectionFilesValue []string
	DetectResultValue   bool
	DetectErrValue      error
	VersionValue        entities.Version
	VersionErrValue     error
	VersionFilesValue   []string
	FilesChangedValue   []string
	FilesChangedErr     error
	WriteErrValue       error
	DepsValue           []entities.Dependency
	DepsErrValue        error
	UpdateErrValue      error
	CommandsValue       []string
}

func (s *LanguageProviderStub) DetectionFiles() []string { return s.DetectionFilesValue }

func (s *LanguageProviderStub) Detect(_ string) (bool, error) {
	return s.DetectResultValue, s.DetectErrValue
}

func (s *LanguageProviderStub) Language() entities.Language { return s.LangValue }

func (s *LanguageProviderStub) VersionFiles() []string { return s.VersionFilesValue }

func (s *LanguageProviderStub) ReadVersion(_ string) (entities.Version, error) {
	return s.VersionValue, s.VersionErrValue
}

func (s *LanguageProviderStub) WriteVersion(_ string, _ entities.Version) error {
	return s.WriteErrValue
}

func (s *LanguageProviderStub) FilesChanged(_ string) ([]string, error) {
	return s.FilesChangedValue, s.FilesChangedErr
}

func (s *LanguageProviderStub) ReadDependencies(_ string) ([]entities.Dependency, error) {
	return s.DepsValue, s.DepsErrValue
}

func (s *LanguageProviderStub) UpdateAll(_ string) error { return s.UpdateErrValue }

func (s *LanguageProviderStub) Commands() []string { return s.CommandsValue }
