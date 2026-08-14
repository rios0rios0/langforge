package javamaven

import (
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/java"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// NewProvider creates a new Java/Maven language provider.
func NewProvider() *repositories.CompositeProvider {
	runner := cmdexec.NewDefaultRunner()
	return &repositories.CompositeProvider{
		LanguageDetector:  &Detector{},
		VersionReader:     &VersionReader{},
		VersionWriter:     &VersionWriter{},
		DependencyReader:  &DependencyReader{},
		DependencyUpdater: NewDependencyUpdater(runner),
		BuildValidator:    NewBuildValidator(runner),
		RuntimeManager:    java.NewRuntimeManager(runner, "mvn spring-boot:run"),
	}
}
