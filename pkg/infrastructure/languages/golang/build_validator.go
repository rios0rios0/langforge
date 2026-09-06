package golang

import (
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// BuildValidator runs lint and build checks for Go projects.
type BuildValidator struct {
	*toolchain.BuildValidator
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{toolchain.NewBuildValidator(runner, goCommands())}
}

// goCommands lists the checks a Go project must pass: golangci-lint, then a
// build of every package.
func goCommands() toolchain.Commands {
	return toolchain.Commands{
		Lint:  []cmdexec.CommandSpec{{Name: "golangci-lint", Args: []string{"run", "./..."}}},
		Build: []cmdexec.CommandSpec{{Name: "go", Args: []string{"build", "./..."}}},
	}
}
