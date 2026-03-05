package csharp

import (
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/support/exec"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// DependencyUpdater runs dotnet outdated --upgrade.
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
		"dotnet outdated --upgrade",
	}
}

// FilesChanged returns the files modified by an update.
func (u *DependencyUpdater) FilesChanged(repoPath string) ([]string, error) {
	matches, err := fileutil.GlobFiles(repoPath, "*.csproj")
	if err != nil {
		return nil, err
	}
	result := make([]string, len(matches))
	for i, m := range matches {
		result[i] = filepath.Join(repoPath, filepath.Base(m))
	}
	return result, nil
}

// UpdateAll runs dotnet outdated --upgrade.
func (u *DependencyUpdater) UpdateAll(repoPath string) error {
	return u.runner.Run(repoPath, "dotnet", "outdated", "--upgrade")
}
