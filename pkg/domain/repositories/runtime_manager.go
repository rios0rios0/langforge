package repositories

// RuntimeManager provides SDK/runtime information and management for a language.
type RuntimeManager interface {
	// SDKName returns the human-readable SDK name (e.g. "Go", "Node.js", "Python").
	SDKName() string

	// VersionManager returns the version manager name (e.g. "gvm", "nvm", "pyenv", "sdkman").
	VersionManager() string

	// InstallCommand returns the shell command to install a specific SDK version.
	InstallCommand(version string) string

	// CurrentVersion returns the currently installed SDK version.
	// Returns empty string and nil error if the SDK is not installed.
	CurrentVersion() (string, error)

	// StartCommand returns the default command to run/start a project.
	StartCommand() string

	// StopCommand returns the default command to stop a project (empty if not applicable).
	StopCommand() string
}
