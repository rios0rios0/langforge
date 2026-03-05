package python

import (
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// DependencyUpdater runs poetry update or pip install -U.
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
		"poetry update",
	}
}

// FilesChanged returns the files modified by an update.
func (u *DependencyUpdater) FilesChanged(repoPath string) ([]string, error) {
	files := []string{}
	if fileutil.Exists(filepath.Join(repoPath, "poetry.lock")) {
		files = append(files, filepath.Join(repoPath, "poetry.lock"))
	}
	if fileutil.Exists(filepath.Join(repoPath, "requirements.txt")) {
		files = append(files, filepath.Join(repoPath, "requirements.txt"))
	}
	return files, nil
}

// UpdateAll runs poetry update if poetry.lock exists, otherwise pip install -U.
func (u *DependencyUpdater) UpdateAll(repoPath string) error {
	if fileutil.Exists(filepath.Join(repoPath, "poetry.lock")) {
		return u.runner.Run(repoPath, "poetry", "update")
	}
	return u.runner.Run(repoPath, "pip", "install", "--upgrade", "-r", "requirements.txt")
}
