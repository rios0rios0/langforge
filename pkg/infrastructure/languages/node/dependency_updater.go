package node

import (
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/support/exec"
)

// DependencyUpdater runs npm update && npm install.
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
		"npm update",
		"npm install",
	}
}

// FilesChanged returns the files modified by an update.
func (u *DependencyUpdater) FilesChanged(repoPath string) ([]string, error) {
	return []string{
		filepath.Join(repoPath, "package.json"),
		filepath.Join(repoPath, "package-lock.json"),
	}, nil
}

// UpdateAll runs npm update && npm install.
func (u *DependencyUpdater) UpdateAll(repoPath string) error {
	if err := u.runner.Run(repoPath, "npm", "update"); err != nil {
		return err
	}
	return u.runner.Run(repoPath, "npm", "install")
}
