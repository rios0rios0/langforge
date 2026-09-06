package csharp

import (
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// BuildValidator runs lint and build checks for C# projects.
type BuildValidator struct {
	*toolchain.BuildValidator
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{toolchain.NewBuildValidator(runner, dotnetCommands())}
}

// dotnetCommands lists the checks a C# project must pass: a formatting check
// that rejects unformatted files, then a build.
func dotnetCommands() toolchain.Commands {
	return toolchain.Commands{
		Lint:  []cmdexec.CommandSpec{{Name: dotnetCLI, Args: []string{"format", "--verify-no-changes"}}},
		Build: []cmdexec.CommandSpec{{Name: dotnetCLI, Args: []string{"build"}}},
	}
}
