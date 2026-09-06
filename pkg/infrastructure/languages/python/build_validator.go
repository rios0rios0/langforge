package python

import (
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// BuildValidator runs lint and build checks for Python projects.
type BuildValidator struct {
	*toolchain.BuildValidator
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{toolchain.NewBuildValidator(runner, pythonCommands())}
}

// pythonCommands lists the checks a Python project must pass: ruff, then a pdm
// build of the package.
func pythonCommands() toolchain.Commands {
	return toolchain.Commands{
		Lint:  []cmdexec.CommandSpec{{Name: "ruff", Args: []string{"check", "."}}},
		Build: []cmdexec.CommandSpec{{Name: "pdm", Args: []string{"build"}}},
	}
}
