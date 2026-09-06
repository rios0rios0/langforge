package javagradle

import (
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// BuildValidator runs lint and build checks for Java/Gradle projects.
type BuildValidator struct {
	*toolchain.BuildValidator
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{toolchain.NewBuildValidator(runner, gradleCommands())}
}

// gradleCommands lists the checks a Gradle project must pass: the wrapper's
// check task, then a build that skips the tests.
func gradleCommands() toolchain.Commands {
	return toolchain.Commands{
		Lint:  []cmdexec.CommandSpec{{Name: "./gradlew", Args: []string{"check"}}},
		Build: []cmdexec.CommandSpec{{Name: "./gradlew", Args: []string{"build", "-x", "test"}}},
	}
}
