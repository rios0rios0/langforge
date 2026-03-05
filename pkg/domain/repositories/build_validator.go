package repositories

// BuildValidator validates that a project builds and passes linting after changes.
type BuildValidator interface {
	// Validate runs build/compile validation commands in the given repo path.
	Validate(repoPath string) error

	// LintCommands returns the ordered list of lint/format commands.
	LintCommands() []string

	// BuildCommands returns the ordered list of build/compile commands.
	BuildCommands() []string
}
