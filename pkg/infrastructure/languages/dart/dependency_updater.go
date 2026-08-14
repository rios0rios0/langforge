package dart

import (
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// DependencyUpdater runs pub's native upgrade toolchain.
type DependencyUpdater struct {
	runner cmdexec.Runner
}

// NewDependencyUpdater creates a DependencyUpdater with the given runner.
func NewDependencyUpdater(runner cmdexec.Runner) *DependencyUpdater {
	return &DependencyUpdater{runner: runner}
}

// Commands returns the shell commands that UpdateAll runs.
//
// The project's toolchain is unknown without a path, so the listing names the
// Dart form; UpdateAll resolves "dart" to "flutter" for a Flutter project.
func (u *DependencyUpdater) Commands() []string {
	return []string{
		"dart pub upgrade --major-versions",
		"dart pub get",
	}
}

// FilesChanged returns the files modified by an update.
func (u *DependencyUpdater) FilesChanged(repoPath string) ([]string, error) {
	files := []string{filepath.Join(repoPath, Pubspec)}
	if lock := filepath.Join(repoPath, PubspecLock); fileutil.Exists(lock) {
		files = append(files, lock)
	}
	return files, nil
}

// UpdateAll runs pub upgrade followed by pub get.
//
// The --major-versions flag is what makes this an upgrade rather than a
// re-resolution: a plain "pub upgrade" only rewrites pubspec.lock within the
// constraints already declared, while --major-versions raises those constraints
// in pubspec.yaml to what pub reports as resolvable. Pub applies that rewrite
// through yaml_edit, so the manifest's comments survive it.
func (u *DependencyUpdater) UpdateAll(repoPath string) error {
	bin := executable(repoPath)
	if err := u.runner.Run(repoPath, bin, "pub", "upgrade", "--major-versions"); err != nil {
		return err
	}
	return u.runner.Run(repoPath, bin, "pub", "get")
}
