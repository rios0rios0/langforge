package toolchain

import (
	"slices"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// Commands lists what a BuildValidator runs, in the order it runs them: every
// lint command first, then every build command.
type Commands struct {
	// Lint holds the lint and format checks.
	Lint []cmdexec.CommandSpec
	// Build holds the build and compile steps, run once every lint check passes.
	Build []cmdexec.CommandSpec
}

// BuildValidator runs a language's lint commands followed by its build commands
// through a Runner, stopping at the first one that fails.
type BuildValidator struct {
	runner   cmdexec.Runner
	commands Commands
}

// NewBuildValidator creates a BuildValidator that runs commands through runner.
func NewBuildValidator(runner cmdexec.Runner, commands Commands) *BuildValidator {
	return &BuildValidator{runner: runner, commands: commands}
}

// LintCommands returns the lint commands as displayable strings.
func (v *BuildValidator) LintCommands() []string {
	return cmdexec.CommandStrings(v.commands.Lint)
}

// BuildCommands returns the build commands as displayable strings.
func (v *BuildValidator) BuildCommands() []string {
	return cmdexec.CommandStrings(v.commands.Build)
}

// Validate runs the lint commands and then the build commands in repoPath and
// returns the first error, leaving the remaining commands unrun.
func (v *BuildValidator) Validate(repoPath string) error {
	for _, cmd := range slices.Concat(v.commands.Lint, v.commands.Build) {
		if err := v.runner.Run(repoPath, cmd.Name, cmd.Args...); err != nil {
			return err
		}
	}
	return nil
}
