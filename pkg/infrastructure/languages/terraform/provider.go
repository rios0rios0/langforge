package terraform

import (
	"github.com/rios0rios0/langforge/pkg/support/exec"
)

// Provider is the composite Terraform language provider.
type Provider struct {
	*Detector
	*VersionReader
	*VersionWriter
	*DependencyReader
	*DependencyUpdater
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

// NewProvider creates a new Terraform language provider.
func NewProvider() *Provider {
	return &Provider{
		Detector:          &Detector{},
		VersionReader:     &VersionReader{},
		VersionWriter:     &VersionWriter{},
		DependencyReader:  &DependencyReader{},
		DependencyUpdater: NewDependencyUpdater(exec.NewDefaultRunner()),
	}
}
