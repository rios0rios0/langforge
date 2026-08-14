package node

import (
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// Provider is the composite Node language provider.
type Provider struct {
	*Detector
	*VersionReader
	*VersionWriter
	*DependencyReader
	*DependencyUpdater
	*BuildValidator
	*RuntimeManager
}

// NewProvider creates a new Node language provider.
func NewProvider() *Provider {
	runner := cmdexec.NewDefaultRunner()
	return &Provider{
		Detector:          &Detector{},
		VersionReader:     &VersionReader{},
		VersionWriter:     &VersionWriter{},
		DependencyReader:  &DependencyReader{},
		DependencyUpdater: NewDependencyUpdater(runner),
		BuildValidator:    NewBuildValidator(runner),
		RuntimeManager:    NewRuntimeManager(runner),
	}
}

// FilesChanged resolves the ambiguity between VersionWriter.FilesChanged and
// DependencyUpdater.FilesChanged by merging both results.
func (p *Provider) FilesChanged(repoPath string) ([]string, error) {
	return repositories.MergeFilesChanged(p.VersionWriter, p.DependencyUpdater, repoPath)
}
