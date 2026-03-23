package python

import "github.com/rios0rios0/langforge/pkg/support/cmdexec"

// BuildValidator runs lint and build checks for Python projects.
type BuildValidator struct {
	runner cmdexec.Runner
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{runner: runner}
}

// LintCommands returns the Python lint commands.
func (v *BuildValidator) LintCommands() []string {
	return []string{"ruff check ."}
}

// BuildCommands returns the Python build commands.
func (v *BuildValidator) BuildCommands() []string {
	return []string{"pdm build"}
}

// Validate runs lint and build commands in the given repo path.
func (v *BuildValidator) Validate(repoPath string) error {
	if err := v.runner.Run(repoPath, "ruff", "check", "."); err != nil {
		return err
	}
	return v.runner.Run(repoPath, "pdm", "build")
}
