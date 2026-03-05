package javagradle

import (
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// DependencyUpdater runs ./gradlew dependencyUpdates.
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
		"./gradlew dependencyUpdates",
	}
}

// FilesChanged returns the files modified by an update.
func (u *DependencyUpdater) FilesChanged(repoPath string) ([]string, error) {
	return []string{
		filepath.Join(repoPath, "build.gradle"),
	}, nil
}

// UpdateAll runs ./gradlew dependencyUpdates.
func (u *DependencyUpdater) UpdateAll(repoPath string) error {
	return u.runner.Run(repoPath, "./gradlew", "dependencyUpdates")
}
