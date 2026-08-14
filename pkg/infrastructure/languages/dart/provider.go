package dart

import (
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// Provider is the composite Dart language provider. It covers Flutter as well:
// both share pubspec.yaml, and the toolchain difference is resolved per
// repository by IsFlutter rather than by a separate provider.
type Provider struct {
	*Detector
	*VersionReader
	*VersionWriter
	*DependencyReader
	*DependencyUpdater
	*BuildValidator
	*RuntimeManager
}

// NewProvider creates a new Dart language provider.
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
