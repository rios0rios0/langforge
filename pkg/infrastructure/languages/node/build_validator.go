package node

import "github.com/rios0rios0/langforge/pkg/support/cmdexec"

// BuildValidator runs lint and build checks for Node.js projects.
type BuildValidator struct {
	runner cmdexec.Runner
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{runner: runner}
}

// LintCommands returns the Node.js lint commands.
func (v *BuildValidator) LintCommands() []string {
	return []string{"npm run lint"}
}

// BuildCommands returns the Node.js build commands.
func (v *BuildValidator) BuildCommands() []string {
	return []string{"npm run build"}
}

// Validate runs lint and build commands in the given repo path.
func (v *BuildValidator) Validate(repoPath string) error {
	if err := v.runner.Run(repoPath, "npm", "run", "lint"); err != nil {
		return err
	}
	return v.runner.Run(repoPath, "npm", "run", "build")
}
