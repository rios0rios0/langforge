package java_gradle

import (
	"github.com/rios0rios0/langforge/pkg/support/exec"
)

// Provider is the composite Java/Gradle language provider.
type Provider struct {
	*Detector
	*VersionReader
	*VersionWriter
	*DependencyReader
	*DependencyUpdater
}

// NewProvider creates a new Java/Gradle language provider.
func NewProvider() *Provider {
	return &Provider{
		Detector:          &Detector{},
		VersionReader:     &VersionReader{},
		VersionWriter:     &VersionWriter{},
		DependencyReader:  &DependencyReader{},
		DependencyUpdater: NewDependencyUpdater(exec.NewDefaultRunner()),
	}
}
