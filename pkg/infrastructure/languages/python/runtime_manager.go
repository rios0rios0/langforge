package python

import (
	"fmt"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// RuntimeManager provides SDK and runtime information for Python projects.
type RuntimeManager struct {
	runner cmdexec.Runner
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{runner: runner}
}

// SDKName returns the human-readable SDK name.
func (m *RuntimeManager) SDKName() string { return "Python" }

// VersionManager returns the version manager name.
func (m *RuntimeManager) VersionManager() string { return "pyenv" }

// StartCommand returns the default command to start a Python project.
func (m *RuntimeManager) StartCommand() string { return "pdm start" }

// StopCommand returns an empty string since there is no standard stop command.
func (m *RuntimeManager) StopCommand() string { return "" }

// InstallCommand returns the pyenv command to install a specific Python version.
func (m *RuntimeManager) InstallCommand(version string) string {
	return fmt.Sprintf("pyenv install %s", version)
}

// CurrentVersion returns the currently installed Python version, or empty if not installed.
func (m *RuntimeManager) CurrentVersion() (string, error) {
	return cmdexec.CapturedVersion(
		m.runner,
		regexp.MustCompile(`Python\s+(\d+\.\d+(?:\.\d+)?)`),
		"python3",
		"--version",
	)
}
