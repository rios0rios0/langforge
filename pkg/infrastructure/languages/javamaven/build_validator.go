package javamaven

import "github.com/rios0rios0/langforge/pkg/support/cmdexec"

// BuildValidator runs lint and build checks for Java/Maven projects.
type BuildValidator struct {
	runner cmdexec.Runner
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{runner: runner}
}

// LintCommands returns the Java/Maven lint commands.
func (v *BuildValidator) LintCommands() []string {
	return []string{"mvn checkstyle:check"}
}

// BuildCommands returns the Java/Maven build commands.
func (v *BuildValidator) BuildCommands() []string {
	return []string{"mvn package -DskipTests"}
}

// Validate runs lint and build commands in the given repo path.
func (v *BuildValidator) Validate(repoPath string) error {
	if err := v.runner.Run(repoPath, "mvn", "checkstyle:check"); err != nil {
		return err
	}
	return v.runner.Run(repoPath, "mvn", "package", "-DskipTests")
}
