package golang

import "github.com/rios0rios0/langforge/pkg/support/cmdexec"

func goLintSpecs() []cmdexec.CommandSpec {
	return []cmdexec.CommandSpec{
		{Name: "golangci-lint", Args: []string{"run", "./..."}},
	}
}

func goBuildSpecs() []cmdexec.CommandSpec {
	return []cmdexec.CommandSpec{
		{Name: "go", Args: []string{"build", "./..."}},
	}
}

// BuildValidator runs lint and build checks for Go projects.
type BuildValidator struct {
	runner cmdexec.Runner
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{runner: runner}
}

// LintCommands returns the Go lint commands.
func (v *BuildValidator) LintCommands() []string {
	return cmdexec.CommandStrings(goLintSpecs())
}

// BuildCommands returns the Go build commands.
func (v *BuildValidator) BuildCommands() []string {
	return cmdexec.CommandStrings(goBuildSpecs())
}

// Validate runs lint and build commands in the given repo path.
func (v *BuildValidator) Validate(repoPath string) error {
	for _, cmd := range append(goLintSpecs(), goBuildSpecs()...) {
		if err := v.runner.Run(repoPath, cmd.Name, cmd.Args...); err != nil {
			return err
		}
	}
	return nil
}
