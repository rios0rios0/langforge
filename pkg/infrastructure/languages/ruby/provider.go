package ruby

import (
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// NewProvider creates a new Ruby language provider.
func NewProvider() *repositories.CompositeProvider {
	runner := cmdexec.NewDefaultRunner()
	return &repositories.CompositeProvider{
		LanguageDetector:  &Detector{},
		VersionReader:     &VersionReader{},
		VersionWriter:     &VersionWriter{},
		DependencyReader:  &DependencyReader{},
		DependencyUpdater: NewDependencyUpdater(runner),
		BuildValidator:    NewBuildValidator(runner),
		RuntimeManager:    NewRuntimeManager(runner),
	}
}
