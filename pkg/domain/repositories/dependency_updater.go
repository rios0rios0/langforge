package repositories

// DependencyUpdater runs the ecosystem's native update toolchain.
type DependencyUpdater interface {
	// UpdateAll runs the ecosystem's native update toolchain.
	UpdateAll(repoPath string) error

	// FilesChanged returns files that may be modified after an update.
	FilesChanged(repoPath string) ([]string, error)

	// Commands returns the ordered list of shell commands that UpdateAll runs,
	// for logging / dry-run support.
	Commands() []string
}
