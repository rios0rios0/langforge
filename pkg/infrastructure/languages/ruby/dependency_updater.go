package ruby

import (
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// DependencyUpdater runs bundle update.
type DependencyUpdater struct {
	runner cmdexec.Runner
}

// NewDependencyUpdater creates a DependencyUpdater with the given runner.
func NewDependencyUpdater(runner cmdexec.Runner) *DependencyUpdater {
	return &DependencyUpdater{runner: runner}
}

// Commands returns the shell commands that UpdateAll runs.
func (u *DependencyUpdater) Commands() []string {
	return []string{
		"bundle update",
	}
}

// FilesChanged returns the files modified by an update.
func (u *DependencyUpdater) FilesChanged(repoPath string) ([]string, error) {
	var files []string
	if fileutil.Exists(filepath.Join(repoPath, "Gemfile.lock")) {
		files = append(files, filepath.Join(repoPath, "Gemfile.lock"))
	}
	return files, nil
}

// UpdateAll runs bundle update.
func (u *DependencyUpdater) UpdateAll(repoPath string) error {
	return u.runner.Run(repoPath, "bundle", "update")
}
