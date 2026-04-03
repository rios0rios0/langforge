package ruby

import "github.com/rios0rios0/langforge/pkg/support/cmdexec"

// BuildValidator runs lint and build checks for Ruby projects.
type BuildValidator struct {
	runner cmdexec.Runner
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{runner: runner}
}

// LintCommands returns the Ruby lint commands.
func (v *BuildValidator) LintCommands() []string {
	return []string{"rubocop"}
}

// BuildCommands returns the Ruby build commands.
func (v *BuildValidator) BuildCommands() []string {
	return []string{"bundle exec rake build"}
}

// Validate runs lint and build commands in the given repo path.
func (v *BuildValidator) Validate(repoPath string) error {
	if err := v.runner.Run(repoPath, "rubocop"); err != nil {
		return err
	}
	return v.runner.Run(repoPath, "bundle", "exec", "rake", "build")
}
