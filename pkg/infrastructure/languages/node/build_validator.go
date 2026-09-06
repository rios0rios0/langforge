package node

import (
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// BuildValidator runs lint and build checks for Node.js projects.
type BuildValidator struct {
	*toolchain.BuildValidator
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{toolchain.NewBuildValidator(runner, npmCommands())}
}

// npmCommands lists the checks a Node.js project must pass: its lint script,
// then its build script.
func npmCommands() toolchain.Commands {
	return toolchain.Commands{
		Lint:  []cmdexec.CommandSpec{{Name: "npm", Args: []string{"run", "lint"}}},
		Build: []cmdexec.CommandSpec{{Name: "npm", Args: []string{"run", "build"}}},
	}
}
