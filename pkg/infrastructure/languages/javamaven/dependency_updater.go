package javamaven

import (
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// DependencyUpdater runs mvn versions:use-latest-releases.
type DependencyUpdater struct {
	runner cmdexec.Runner
}

// NewDependencyUpdater creates a DependencyUpdater with the default runner.
func NewDependencyUpdater(runner cmdexec.Runner) *DependencyUpdater {
	return &DependencyUpdater{runner: runner}
}

// Commands returns the shell commands that UpdateAll runs.
func (u *DependencyUpdater) Commands() []string {
	return []string{
		"mvn versions:use-latest-releases",
	}
}

// FilesChanged returns the files modified by an update.
func (u *DependencyUpdater) FilesChanged(repoPath string) ([]string, error) {
	return []string{
		filepath.Join(repoPath, "pom.xml"),
	}, nil
}

// UpdateAll runs mvn versions:use-latest-releases.
func (u *DependencyUpdater) UpdateAll(repoPath string) error {
	return u.runner.Run(repoPath, "mvn", "versions:use-latest-releases")
}
