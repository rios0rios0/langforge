package node

import (
	"fmt"
	"strings"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// RuntimeManager provides SDK and runtime information for Node.js projects.
type RuntimeManager struct {
	runner cmdexec.Runner
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{runner: runner}
}

// SDKName returns the human-readable SDK name.
func (m *RuntimeManager) SDKName() string { return "Node.js" }

// VersionManager returns the version manager name.
func (m *RuntimeManager) VersionManager() string { return "nvm" }

// StartCommand returns the default command to start a Node.js project.
func (m *RuntimeManager) StartCommand() string { return "npm start" }

// StopCommand returns an empty string since there is no standard stop command.
func (m *RuntimeManager) StopCommand() string { return "" }

// InstallCommand returns the nvm command to install a specific Node.js version.
func (m *RuntimeManager) InstallCommand(version string) string {
	return fmt.Sprintf("nvm install %s", version)
}

// CurrentVersion returns the currently installed Node.js version, or empty if not installed.
func (m *RuntimeManager) CurrentVersion() (string, error) {
	output, err := m.runner.RunOutput(".", "node", "--version")
	if err != nil {
		return "", nil
	}
	return strings.TrimPrefix(output, "v"), nil
}
