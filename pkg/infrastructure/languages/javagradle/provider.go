package javagradle

import (
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
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
		DependencyUpdater: NewDependencyUpdater(cmdexec.NewDefaultRunner()),
	}
}

// FilesChanged resolves the ambiguity between VersionWriter.FilesChanged and
// DependencyUpdater.FilesChanged by merging both results.
func (p *Provider) FilesChanged(repoPath string) ([]string, error) {
	vFiles, err := p.VersionWriter.FilesChanged(repoPath)
	if err != nil {
		return nil, err
	}
	dFiles, err := p.DependencyUpdater.FilesChanged(repoPath)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{}, len(vFiles))
	for _, f := range vFiles {
		seen[f] = struct{}{}
	}
	for _, f := range dFiles {
		if _, ok := seen[f]; !ok {
			vFiles = append(vFiles, f)
		}
	}
	return vFiles, nil
}
