package dart

import "github.com/rios0rios0/langforge/pkg/support/cmdexec"

// BuildValidator runs lint and build checks for Dart and Flutter projects.
type BuildValidator struct {
	runner cmdexec.Runner
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{runner: runner}
}

// LintCommands returns the Dart lint commands.
func (v *BuildValidator) LintCommands() []string {
	return []string{"dart analyze"}
}

// BuildCommands returns the Dart build commands.
//
// "pub get" stands in for a compile step because Dart has no single build
// output: a package produces none, an application targets web, Android, iOS,
// desktop or all of them. Resolving the dependency graph is the check that is
// meaningful for every project shape.
func (v *BuildValidator) BuildCommands() []string {
	return []string{"dart pub get"}
}

// Validate runs analysis and dependency resolution in the given repo path.
func (v *BuildValidator) Validate(repoPath string) error {
	bin := executable(repoPath)
	if err := v.runner.Run(repoPath, bin, "pub", "get"); err != nil {
		return err
	}
	return v.runner.Run(repoPath, bin, "analyze")
}
