package csharp

import "github.com/rios0rios0/langforge/pkg/support/cmdexec"

// BuildValidator runs lint and build checks for C# projects.
type BuildValidator struct {
	runner cmdexec.Runner
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{runner: runner}
}

// LintCommands returns the C# lint commands.
func (v *BuildValidator) LintCommands() []string {
	return []string{"dotnet format --verify-no-changes"}
}

// BuildCommands returns the C# build commands.
func (v *BuildValidator) BuildCommands() []string {
	return []string{"dotnet build"}
}

// Validate runs lint and build commands in the given repo path.
func (v *BuildValidator) Validate(repoPath string) error {
	if err := v.runner.Run(repoPath, "dotnet", "format", "--verify-no-changes"); err != nil {
		return err
	}
	return v.runner.Run(repoPath, "dotnet", "build")
}
