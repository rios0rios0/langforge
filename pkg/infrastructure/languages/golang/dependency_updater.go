package golang

import (
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/support/exec"
)

// DependencyUpdater runs go get -u all && go mod tidy.
type DependencyUpdater struct {
	runner exec.Runner
}

// NewDependencyUpdater creates a DependencyUpdater with the default runner.
func NewDependencyUpdater(runner exec.Runner) *DependencyUpdater {
	return &DependencyUpdater{runner: runner}
}

// Commands returns the shell commands that UpdateAll runs.
func (u *DependencyUpdater) Commands() []string {
	return []string{
		"go get -u all",
		"go mod tidy",
	}
}

// FilesChanged returns the files modified by an update.
func (u *DependencyUpdater) FilesChanged(repoPath string) ([]string, error) {
	return []string{
		filepath.Join(repoPath, "go.mod"),
		filepath.Join(repoPath, "go.sum"),
	}, nil
}

// UpdateAll runs go get -u all && go mod tidy.
func (u *DependencyUpdater) UpdateAll(repoPath string) error {
	if err := u.runner.Run(repoPath, "go", "get", "-u", "all"); err != nil {
		return err
	}
	return u.runner.Run(repoPath, "go", "mod", "tidy")
}
