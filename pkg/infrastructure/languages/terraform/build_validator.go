package terraform

import "github.com/rios0rios0/langforge/pkg/support/cmdexec"

// BuildValidator runs lint and build checks for Terraform projects.
type BuildValidator struct {
	runner cmdexec.Runner
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{runner: runner}
}

// LintCommands returns the Terraform lint commands.
func (v *BuildValidator) LintCommands() []string {
	return []string{"terraform fmt -check"}
}

// BuildCommands returns the Terraform build commands.
func (v *BuildValidator) BuildCommands() []string {
	return []string{"terraform validate"}
}

// Validate runs lint and build commands in the given repo path.
func (v *BuildValidator) Validate(repoPath string) error {
	if err := v.runner.Run(repoPath, "terraform", "fmt", "-check"); err != nil {
		return err
	}
	return v.runner.Run(repoPath, "terraform", "validate")
}
