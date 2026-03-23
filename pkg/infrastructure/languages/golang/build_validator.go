package golang

import "github.com/rios0rios0/langforge/pkg/support/cmdexec"

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
	return []string{"golangci-lint run ./..."}
}

// BuildCommands returns the Go build commands.
func (v *BuildValidator) BuildCommands() []string {
	return []string{"go build ./..."}
}

// Validate runs lint and build commands in the given repo path.
func (v *BuildValidator) Validate(repoPath string) error {
	if err := v.runner.Run(repoPath, "golangci-lint", "run", "./..."); err != nil {
		return err
	}
	return v.runner.Run(repoPath, "go", "build", "./...")
}
