package javamaven

import (
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// BuildValidator runs lint and build checks for Java/Maven projects.
type BuildValidator struct {
	*toolchain.BuildValidator
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{toolchain.NewBuildValidator(runner, mavenCommands())}
}

// mavenCommands lists the checks a Maven project must pass: checkstyle, then a
// package build that skips the tests.
func mavenCommands() toolchain.Commands {
	return toolchain.Commands{
		Lint:  []cmdexec.CommandSpec{{Name: "mvn", Args: []string{"checkstyle:check"}}},
		Build: []cmdexec.CommandSpec{{Name: "mvn", Args: []string{"package", "-DskipTests"}}},
	}
}
