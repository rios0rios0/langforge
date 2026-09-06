package ruby

import (
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// BuildValidator runs lint and build checks for Ruby projects.
type BuildValidator struct {
	*toolchain.BuildValidator
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{toolchain.NewBuildValidator(runner, bundlerCommands())}
}

// bundlerCommands lists the checks a Ruby project must pass: rubocop, then the
// gem's rake build task.
func bundlerCommands() toolchain.Commands {
	return toolchain.Commands{
		Lint:  []cmdexec.CommandSpec{{Name: "rubocop"}},
		Build: []cmdexec.CommandSpec{{Name: "bundle", Args: []string{"exec", "rake", "build"}}},
	}
}
