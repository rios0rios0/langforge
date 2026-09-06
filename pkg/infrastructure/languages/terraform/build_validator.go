package terraform

import (
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// BuildValidator runs lint and build checks for Terraform projects.
type BuildValidator struct {
	*toolchain.BuildValidator
}

// NewBuildValidator creates a BuildValidator with the given runner.
func NewBuildValidator(runner cmdexec.Runner) *BuildValidator {
	return &BuildValidator{toolchain.NewBuildValidator(runner, terraformCommands())}
}

// terraformCommands lists the checks a Terraform project must pass: a
// formatting check, then validation of the configuration.
func terraformCommands() toolchain.Commands {
	return toolchain.Commands{
		Lint:  []cmdexec.CommandSpec{{Name: terraformCLI, Args: []string{"fmt", "-check"}}},
		Build: []cmdexec.CommandSpec{{Name: terraformCLI, Args: []string{"validate"}}},
	}
}
